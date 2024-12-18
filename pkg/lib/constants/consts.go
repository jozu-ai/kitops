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
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

type ConfigKey struct{}

const (
	// Default name for Kitfile (otherwise specified via the -f flag in pack)
	DefaultKitfileName = "Kitfile"
	// IgnoreFileName is the name for the Kit ignore file
	IgnoreFileName = ".kitignore"

	// Constants for the directory structure of kit's cached images and credentials
	// Modelkits are stored in $KITOPS_HOME/storage/ and
	// credentials are stored in $KITOPS_HOME/credentials.json
	DefaultConfigSubdir               = "kitops"
	StorageSubpath                    = "storage"
	CredentialsSubpath                = "credentials.json"
	HarnessSubpath                    = "harness"
	HarnessProcessFile                = "process.pid"
	HarnessLogFile                    = "harness.log"
	UpdateNotificationsConfigFilename = "disable-update-notifications"

	// Kitops-specific annotations for modelkit artifacts
	CliVersionAnnotation = "ml.kitops.modelkit.cli-version"

	// MaxModelRefChain is the maximum number of "parent" modelkits a modelkit may have
	// by e.g. referring to another modelkit in its .model.path
	MaxModelRefChain = 10
)

var (
	localIndexNameRegexp = regexp.MustCompile(`^([-A-Za-z0-9_-]*={0,3})-index.json$`)
)

func DefaultKitfileNames() []string {
	return []string{"Kitfile", "kitfile", ".kitfile"}
}

func IsDefaultKitfileName(filename string) bool {
	for _, name := range DefaultKitfileNames() {
		if name == filename {
			return true
		}
	}
	return false
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

func IngestPath(storageBase string) string {
	return filepath.Join(storageBase, "ingest")
}

func HarnessPath(configBase string) string {
	return filepath.Join(configBase, HarnessSubpath)
}

func CredentialsPath(configBase string) string {
	return filepath.Join(configBase, CredentialsSubpath)
}

// IndexJsonPath is a wrapper for getting the index.json path for a local OCI index,
// based off the base path of the index.
func IndexJsonPath(storageBase string) string {
	return filepath.Join(storageBase, "index.json")
}

// IndexJsonPathForRepo returns the path to an index.json that is scoped for a specific repo (org/name)
func IndexJsonPathForRepo(storageBase, repo string) string {
	// We need to encode the repo as it may contain invalid characters for filenames
	repoEncoded := base64.URLEncoding.EncodeToString([]byte(repo))
	indexFileName := fmt.Sprintf("%s-index.json", repoEncoded)
	return filepath.Join(storageBase, indexFileName)
}

// RepoForIndexJsonPath returns the repository for an OCI index JSON file as generated
// by IndexJsonPathForRepo
func RepoForIndexJsonPath(indexPath string) (string, error) {
	filename := filepath.Base(indexPath)
	matches := localIndexNameRegexp.FindStringSubmatch(filename)
	if len(matches) == 0 {
		return "", fmt.Errorf("invalid local OCI index name: %s", filename)
	}

	repoBytes, err := base64.URLEncoding.DecodeString(matches[1])
	if err != nil {
		return "", fmt.Errorf("failed to parse repo from index %s: %w", indexPath, err)
	}
	return string(repoBytes), nil
}

func FileIsLocalIndex(indexPath string) bool {
	filename := filepath.Base(indexPath)
	return localIndexNameRegexp.MatchString(filename)
}

func TagIndexPathForRepo(storageBase, repo string) string {
	// We need to encode the repo as it may contain invalid characters for filenames
	repoEncoded := base64.URLEncoding.EncodeToString([]byte(repo))
	indexFileName := fmt.Sprintf("%s-tags.json", repoEncoded)
	return filepath.Join(storageBase, indexFileName)
}
