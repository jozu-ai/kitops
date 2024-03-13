package repo

import (
	"context"
	"fmt"
	"kitops/pkg/lib/network"

	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

type RegistryOptions struct {
	PlainHTTP       bool
	SkipTLSVerify   bool
	CredentialsPath string
}

// NewRegistry returns a new *remote.Registry for hostname, with credentials and TLS
// configured.
func NewRegistry(hostname string, opts *RegistryOptions) (*remote.Registry, error) {
	reg, err := remote.NewRegistry(hostname)
	if err != nil {
		return nil, err
	}

	reg.PlainHTTP = opts.PlainHTTP
	credentialStore, err := network.NewCredentialStore(opts.CredentialsPath)
	if err != nil {
		return nil, err
	}
	client := network.ClientWithAuth(credentialStore, &network.ClientOpts{TLSSkipVerify: opts.SkipTLSVerify})
	reg.Client = client

	return reg, nil
}

func NewRepository(ctx context.Context, hostname, repository string, opts *RegistryOptions) (registry.Repository, error) {
	registry, err := NewRegistry(hostname, opts)
	if err != nil {
		return nil, err
	}
	repo, err := registry.Repository(ctx, repository)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}
	return repo, nil
}
