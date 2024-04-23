/*
Copyright © 2024 Jozu.com
*/
package pull

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"kitops/pkg/lib/constants"
	kfutils "kitops/pkg/lib/kitfile"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

func runPull(ctx context.Context, opts *pullOptions) (ocispec.Descriptor, error) {
	storageHome := constants.StoragePath(opts.configHome)
	localStorePath := repo.RepoPath(storageHome, opts.modelRef)
	localStore, err := oci.New(localStorePath)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	return runPullRecursive(ctx, localStore, opts, []string{})
}

func runPullRecursive(ctx context.Context, localStore *oci.Store, opts *pullOptions, pulledRefs []string) (ocispec.Descriptor, error) {
	refStr := repo.FormatRepositoryForDisplay(opts.modelRef.String())
	if idx := getIndex(pulledRefs, refStr); idx != -1 {
		cycleStr := fmt.Sprintf("[%s=>%s]", strings.Join(pulledRefs[idx:], "=>"), refStr)
		return ocispec.DescriptorEmptyJSON, fmt.Errorf("Found cycle in modelkit references: %s", cycleStr)
	}
	pulledRefs = append(pulledRefs, refStr)
	if len(pulledRefs) > constants.MaxModelRefChain {
		return ocispec.DescriptorEmptyJSON, fmt.Errorf("Reached maximum number of model references: [%s]", strings.Join(pulledRefs, "=>"))
	}

	remoteRegistry, err := repo.NewRegistry(opts.modelRef.Registry, &repo.RegistryOptions{
		PlainHTTP:       opts.PlainHTTP,
		SkipTLSVerify:   !opts.TlsVerify,
		CredentialsPath: constants.CredentialsPath(opts.configHome),
	})
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}

	desc, err := pullModel(ctx, remoteRegistry, localStore, opts.modelRef)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}

	if err := pullParents(ctx, localStore, desc, opts, pulledRefs); err != nil {
		return ocispec.DescriptorEmptyJSON, fmt.Errorf("Failed to pull referenced modelkits: %w", err)
	}

	return desc, nil
}

func pullParents(ctx context.Context, localStore *oci.Store, desc ocispec.Descriptor, optsIn *pullOptions, pulledRefs []string) error {
	_, config, err := repo.GetManifestAndConfig(ctx, localStore, desc)
	if err != nil {
		return err
	}
	if config.Model == nil || !kfutils.IsModelKitReference(config.Model.Path) {
		return nil
	}
	output.Infof("Pulling referenced image %s", config.Model.Path)
	parentRef, _, err := repo.ParseReference(config.Model.Path)
	if err != nil {
		return err
	}
	opts := *optsIn
	opts.modelRef = parentRef
	_, err = runPullRecursive(ctx, localStore, &opts, pulledRefs)
	return err
}

func pullModel(ctx context.Context, remoteRegistry *remote.Registry, localStore *oci.Store, ref *registry.Reference) (ocispec.Descriptor, error) {
	repo, err := remoteRegistry.Repository(ctx, ref.Repository)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to read repository: %w", err)
	}
	if err := referenceIsModel(ctx, ref, repo); err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	trackedRepo, logger := output.WrapTarget(localStore)
	desc, err := oras.Copy(ctx, repo, ref.Reference, trackedRepo, ref.Reference, oras.DefaultCopyOptions)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to copy to remote: %w", err)
	}
	logger.Wait()

	return desc, err
}

func referenceIsModel(ctx context.Context, ref *registry.Reference, repo registry.Repository) error {
	desc, rc, err := repo.FetchReference(ctx, ref.Reference)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %w", ref.String(), err)
	}
	defer rc.Close()

	if desc.MediaType != ocispec.MediaTypeImageManifest {
		return fmt.Errorf("reference %s is not an image manifest", ref.String())
	}
	manifestBytes, err := io.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}
	manifest := &ocispec.Manifest{}
	if err := json.Unmarshal(manifestBytes, manifest); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}
	if manifest.Config.MediaType != constants.ModelConfigMediaType {
		return fmt.Errorf("reference %s does not refer to a model", ref.String())
	}
	return nil
}

func getIndex(list []string, s string) int {
	for idx, item := range list {
		if s == item {
			return idx
		}
	}
	return -1
}
