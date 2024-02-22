/*
Copyright © 2024 Jozu.com
*/
package pull

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kitops/pkg/lib/constants"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

func pullModel(ctx context.Context, remoteRegistry *remote.Registry, localStore *oci.Store, ref *registry.Reference) (ocispec.Descriptor, error) {
	repo, err := remoteRegistry.Repository(ctx, ref.Repository)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to read repository: %w", err)
	}
	if err := referenceIsModel(ctx, ref, repo); err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}

	desc, err := oras.Copy(ctx, repo, ref.Reference, localStore, ref.Reference, oras.DefaultCopyOptions)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to copy to remote: %w", err)
	}

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
