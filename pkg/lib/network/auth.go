package network

import (
	"kitops/pkg/cmd/version"
	"net/http"

	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/credentials"
	"oras.land/oras-go/v2/registry/remote/retry"
)

type ClientOpts struct {
	TLSSkipVerify bool
}

func NewCredentialStore(storePath string) (credentials.Store, error) {
	return credentials.NewStore(storePath, credentials.StoreOptions{
		DetectDefaultNativeStore: true,
		AllowPlaintextPut:        true,
	})
}

func ClientWithAuth(store credentials.Store, opts *ClientOpts) *auth.Client {
	client := DefaultClient(opts)
	client.Credential = credentials.Credential(store)

	return client
}

func DefaultClient(opts *ClientOpts) *auth.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if opts.TLSSkipVerify {
		transport.TLSClientConfig.InsecureSkipVerify = true
	}

	client := &auth.Client{
		Client: &http.Client{
			Transport: retry.NewTransport(transport),
		},
		Cache: auth.NewCache(),
		Header: http.Header{
			"User-Agent": {"kitops-cli/" + version.Version},
		},
	}

	return client
}
