package info

import (
	"context"
	"fmt"
	"kitops/pkg/artifact"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"
)

func getInfo(ctx context.Context, opts *infoOptions) (*artifact.KitFile, error) {
	if opts.checkRemote {
		return getRemoteConfig(ctx, opts)
	} else {
		return getLocalConfig(ctx, opts)
	}
}

func getLocalConfig(ctx context.Context, opts *infoOptions) (*artifact.KitFile, error) {
	storageRoot := constants.StoragePath(opts.configHome)
	store, err := repo.NewLocalStore(storageRoot, opts.modelRef)
	if err != nil {
		return nil, fmt.Errorf("failed to read local storage: %w", err)
	}
	_, config, err := repo.ResolveManifestAndConfig(ctx, store, opts.modelRef.Reference)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getRemoteConfig(ctx context.Context, opts *infoOptions) (*artifact.KitFile, error) {
	repository, err := repo.NewRepository(ctx, opts.modelRef.Registry, opts.modelRef.Repository, &repo.RegistryOptions{
		PlainHTTP:       opts.PlainHTTP,
		SkipTLSVerify:   !opts.TlsVerify,
		CredentialsPath: constants.CredentialsPath(opts.configHome),
	})
	if err != nil {
		return nil, err
	}
	_, config, err := repo.ResolveManifestAndConfig(ctx, repository, opts.modelRef.Reference)
	if err != nil {
		return nil, err
	}
	return config, nil
}
