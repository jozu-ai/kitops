package inspect

import (
	"context"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func inspectReference(ctx context.Context, opts *inspectOptions) (*ocispec.Manifest, error) {
	if opts.checkRemote {
		return getRemoteManifest(ctx, opts)
	} else {
		return getLocalManifest(ctx, opts)
	}
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
