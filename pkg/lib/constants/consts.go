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

package constants

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type ConfigKey struct{}

const (
	// Default name for Kitfile (otherwise specified via the -f flag in pack)
	DefaultKitfileName = "Kitfile"
	// IgnoreFileName is the name for the Kit ignore file
	IgnoreFileName = ".kitignore"

	// Constants for the directory structure of kit's cached images and credentials
	// Modelkits are stored in <configpath>/kitops/storage/ and
	// credentials are stored in <configpath>/kitops/credentials.json
	DefaultConfigSubdir = "kitops"
	StorageSubpath      = "storage"
	CredentialsSubpath  = "credentials.json"

	// Media type for the model layer
	ModelLayerMediaType = "application/vnd.kitops.modelkit.model.v1.tar+gzip"
	// Media type for model part layer
	ModelPartLayerMediaType = "application/vnd.kitops.modelkit.modelpart.v1.tar+gzip"
	// Media type for the dataset layer
	DataSetLayerMediaType = "application/vnd.kitops.modelkit.dataset.v1.tar+gzip"
	// Media type for the code layer
	CodeLayerMediaType = "application/vnd.kitops.modelkit.code.v1.tar+gzip"
	// Media type for the model config (Kitfile)
	ModelConfigMediaType = "application/vnd.kitops.modelkit.config.v1+json"

	// MaxModelRefChain is the maximum number of "parent" modelkits a modelkit may have
	// by e.g. referring to another modelkit in its .model.path
	MaxModelRefChain = 10
)

func DefaultKitfileNames() []string {
	return []string{"Kitfile", "kitfile", ".kitfile"}
}

// DefaultConfigPath returns the default configuration and cache directory for the CLI.
// This is platform-dependent, using
//   - $XDG_DATA_HOME/kitops on Linux, with fall back to $HOME/.local/share/kitops
//   - ~/Library/Caches/kitops on MacOS
//   - %LOCALAPPDATA%\kitops
func DefaultConfigPath() (string, error) {
	switch runtime.GOOS {
	case "linux":
		datahome := os.Getenv("XDG_DATA_HOME")
		if datahome == "" {
			// Use default ~/.local/share/
			userhome := os.Getenv("HOME")
			if userhome == "" {
				return "", fmt.Errorf("could not get $HOME directory")
			}
			datahome = filepath.Join(userhome, ".local", "share")
		}
		return filepath.Join(datahome, DefaultConfigSubdir), nil

	case "darwin":
		// Use ~/Library/Caches/kitops
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(cacheDir, DefaultConfigSubdir), nil

	case "windows":
		// Use %LOCALAPPDATA%\kitops
		appdata, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(appdata, DefaultConfigSubdir), nil

	default:
		return "", fmt.Errorf("Unrecognized operating system")
	}
}

func StoragePath(configBase string) string {
	return filepath.Join(configBase, StorageSubpath)
}

func CredentialsPath(configBase string) string {
	return filepath.Join(configBase, CredentialsSubpath)
}

// IndexJsonPath is a wrapper for getting the index.json path for a local OCI index,
// based off the base path of the index.
func IndexJsonPath(storageBase string) string {
	return filepath.Join(storageBase, "index.json")
}
