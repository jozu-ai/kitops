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

package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type remoteErr struct {
	Method       string
	URL          *url.URL
	StatusCode   int
	Errors       []jsonError `json:"errors,omitempty"`
	ExtraMessage string
}

type jsonError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  any    `json:"detail,omitempty"`
}

func (e remoteErr) Error() string {
	var errMsg string

	getText := func(je jsonError) string {
		if je.Message == "" {
			return http.StatusText(e.StatusCode)
		}
		return fmt.Sprintf("%s: %s", strings.ToLower(je.Code), je.Message)
	}

	switch len(e.Errors) {
	case 0:
		errMsg = http.StatusText(e.StatusCode)
	case 1:
		errMsg = getText(e.Errors[0])
	default:
		var msgs []string
		for _, err := range e.Errors {
			msgs = append(msgs, getText(err))
		}
		errMsg = fmt.Sprintf("multiple errors: %s", strings.Join(msgs, "; "))
	}
	if e.ExtraMessage != "" {
		errMsg = fmt.Sprintf("%s (additional info: %s)", errMsg, e.ExtraMessage)
	}
	return errMsg
}

func handleRemoteError(resp *http.Response) error {
	respErr := &remoteErr{
		Method:     resp.Request.Method,
		URL:        resp.Request.URL,
		StatusCode: resp.StatusCode,
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8*1024))
	if err != nil {
		respErr.ExtraMessage = fmt.Sprintf("failed to read response body: %s", err)
		return respErr
	}

	if resp.Header.Get("Content-Type") == "application/json" {
		if err := json.Unmarshal(body, &respErr); err != nil {
			respErr.ExtraMessage = fmt.Sprintf("failed to unmarshal response body: %s, body: %s", err, string(body))
			return respErr
		}
	} else if len(body) > 0 {
		respErr.ExtraMessage = string(body)
	}

	return respErr
}
