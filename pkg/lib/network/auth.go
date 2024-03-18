// Copyright 2024 The KitOps Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

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

// ClientWithAuth returns a default *auth.Client using the provided credentials
// store
func ClientWithAuth(store credentials.Store, opts *ClientOpts) *auth.Client {
	client := DefaultClient(opts)
	client.Credential = credentials.Credential(store)

	return client
}

// DefaultClient returns an *auth.Client with a default User-Agent header and TLS
// configured from opts (optionally disabling TLS verification)
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
