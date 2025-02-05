// Copyright 2025 The KitOps Authors.
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

package hf

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"kitops/pkg/output"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

const (
	resolveURLFmt = "https://huggingface.co/%s/resolve/main/%s"
)

func DownloadFiles(ctx context.Context, modelRepo, destDir string, filepaths []string, token string) error {
	client := &http.Client{}

	sem := semaphore.NewWeighted(5)
	errs, errCtx := errgroup.WithContext(ctx)
	var semErr error

	progress, plog := output.NewDownloadProgress()

	for _, f := range filepaths {
		f := f
		if err := sem.Acquire(errCtx, 1); err != nil {
			semErr = err
			break
		}

		fileURL := fmt.Sprintf(resolveURLFmt, modelRepo, f)
		destPath := filepath.Join(destDir, f)
		errs.Go(func() error {
			defer sem.Release(1)
			plog.Infof("Downloading file %s", f)
			return downloadFile(errCtx, client, token, fileURL, destPath, f, progress, plog)
		})
	}

	if err := errs.Wait(); err != nil {
		return err
	}
	if semErr != nil {
		return fmt.Errorf("failed to acquire lock: %w", semErr)
	}
	progress.Done()

	return nil
}

func downloadFile(
	ctx context.Context,
	client *http.Client,
	token, srcURL, destPath, filename string,
	progress *output.DownloadProgressBar,
	plog *output.ProgressLogger) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srcURL, nil)
	if err != nil {
		return fmt.Errorf("failed to resolve URL: %w", err)
	}
	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling API: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			plog.Logf(output.LogLevelWarn, "Failed to close response body: %w", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received status code %d when downloading file %s", resp.StatusCode, destPath)
	}

	contentRC := progress.TrackDownload(resp.Body, filename, resp.ContentLength)
	defer func() {
		if err := contentRC.Close(); err != nil {
			plog.Logf(output.LogLevelWarn, "TEMP: see if this is an issue: %s", err)
		}
	}()

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.OpenFile(destPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil && !errors.Is(err, fs.ErrClosed) {
			plog.Logf(output.LogLevelError, "Error closing file %s: %s", destPath, err)
		}
	}()

	n, err := io.Copy(f, contentRC)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	if resp.ContentLength > 0 && n != resp.ContentLength {
		return fmt.Errorf("mismatched file size: expected %d but got %d", resp.ContentLength, n)
	}

	return nil
}
