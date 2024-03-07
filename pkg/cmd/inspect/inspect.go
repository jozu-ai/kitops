package inspect

import (
	"context"
	"errors"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/errdef"
)

func inspectReference(ctx context.Context, opts *inspectOptions) (*ocispec.Manifest, error) {
	if opts.modelRef.Registry == repo.DefaultRegistry {
		// Local only check
		return getLocalManifest(ctx, opts)
	}
	if opts.checkRemote {
		// Remote only check
		return getRemoteManifest(ctx, opts)
	}

	// Check locally first; if not found check remote
	manifest, err := getLocalManifest(ctx, opts)
	if err == nil {
		return manifest, nil
	} else if !errors.Is(err, errdef.ErrNotFound) {
		return nil, err
	}
	output.Debugf("ModelKit not found locally, checking remote.")
	return getRemoteManifest(ctx, opts)
}

func getLocalManifest(ctx context.Context, opts *inspectOptions) (*ocispec.Manifest, error) {
	storageRoot := constants.StoragePath(opts.configHome)
	store, err := repo.NewLocalStore(storageRoot, opts.modelRef)
	if err != nil {
		return nil, fmt.Errorf("failed to read local storage: %w", err)
	}
	return repo.ResolveManifest(ctx, store, opts.modelRef.Reference)
}

func getRemoteManifest(ctx context.Context, opts *inspectOptions) (*ocispec.Manifest, error) {
	repository, err := repo.NewRepository(ctx, opts.modelRef.Registry, opts.modelRef.Repository, &repo.RegistryOptions{
		PlainHTTP:       opts.PlainHTTP,
		SkipTLSVerify:   !opts.TlsVerify,
		CredentialsPath: constants.CredentialsPath(opts.configHome),
	})
	if err != nil {
		return nil, err
	}
	return repo.ResolveManifest(ctx, repository, opts.modelRef.Reference)
}
