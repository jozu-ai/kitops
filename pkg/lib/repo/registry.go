package repo

import (
	"kitops/pkg/lib/network"

	"oras.land/oras-go/v2/registry/remote"
)

type RegistryOptions struct {
	PlainHTTP       bool
	SkipTLSVerify   bool
	CredentialsPath string
}

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
