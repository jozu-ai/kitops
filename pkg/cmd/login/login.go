package login

import (
	"context"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/network"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"

	"oras.land/oras-go/v2/registry/remote/credentials"
)

func login(ctx context.Context, opts *loginOptions) error {
	credentialsStorePath := constants.CredentialsPath(opts.configHome)
	store, err := network.NewCredentialStore(credentialsStorePath)
	if err != nil {
		return err
	}
	registry, err := repo.NewRegistry(opts.registry, &repo.RegistryOptions{
		PlainHTTP:       opts.PlainHTTP,
		SkipTLSVerify:   !opts.TlsVerify,
		CredentialsPath: credentialsStorePath,
	})
	if err != nil {
		return fmt.Errorf("could not resolve registry %s: %w", opts.registry, err)
	}
	if err := credentials.Login(ctx, store, registry, opts.credential); err != nil {
		return err
	}
	output.Infoln("Log in successful")
	return nil
}
