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

package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"kitops/pkg/lib/constants"
	"kitops/pkg/output"

	"golang.org/x/mod/semver"
)

const releaseUrl = "https://api.github.com/repos/jozu-ai/kitops/releases/latest"

// Regexp for a semver version -- taken from https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
// We've added an optional 'v' to the start (e.g. v1.2.3) since using a 'v' prefix is common (and used, in our case)
// Capture groups are:
//
//	[1] - Major version
//	[2] - Minor version
//	[3] - Bugfix/z-stream version
//	[4] - Pre-release identifiers (1.2.3-<info>), if present
//	[5] - Build metadata (1.2.3+<metadata>), if present
var versionTagRegexp = regexp.MustCompile(`^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

type ghReleaseInfo struct {
	TagName    string `json:"tag_name"`
	Prerelease bool   `json:"prerelease"`
	Draft      bool   `json:"draft"`
	Url        string `json:"html_url"`
}

func CheckForUpdate() {
	// If this isn't a release version of kit, don't nag the user unnecessarily
	if constants.Version == "unknown" || !versionTagRegexp.MatchString(constants.Version) {
		return
	}

	resp, err := http.Get(releaseUrl)
	if err != nil {
		output.Debugf("Failed to check for updates: %s", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		output.Debugf("Failed to read GitHub response body: %s", err)
		return
	}
	info := &ghReleaseInfo{}
	if err := json.Unmarshal(respBody, info); err != nil {
		output.Debugf("Failed to parse GitHub response body: %s", err)
		return
	}
	if info.Prerelease || info.Draft {
		// This isn't a full release; for now just don't notify users, even if there is a newer full release we don't know about
		return
	}

	// The Go semver package requires versions start with a 'v' (contrary to the spec)
	currentVersion := fmt.Sprintf("v%s", strings.TrimPrefix(constants.Version, "v"))
	latestVersion := fmt.Sprintf("v%s", strings.TrimPrefix(info.TagName, "v"))
	if semver.Compare(currentVersion, latestVersion) < 0 {
		output.Infof("Note: A new version of Kit is available! You are using Kit %s. The latest version is %s.", currentVersion, latestVersion)
		output.Infof("      To see a list of changes, visit %s", info.Url)
		output.Infof("") // Add a newline to not confuse it with regular output
	}
}
