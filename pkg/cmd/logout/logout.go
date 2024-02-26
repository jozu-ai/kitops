package logout

import (
	"context"
	"kitops/pkg/lib/network"
	"kitops/pkg/output"

	"oras.land/oras-go/v2/registry/remote/credentials"
)

func logout(ctx context.Context, hostname string, credentialsPath string) error {
	store, err := network.NewCredentialStore(credentialsPath)
	if err != nil {
		return err
	}
	if err := credentials.Logout(ctx, store, hostname); err != nil {
		return err
	}
	output.Infof("Successfully logged out from %s", hostname)
	return nil
}
