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
