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

package output

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"golang.org/x/term"
	"oras.land/oras-go/v2"
)

func shouldPrintProgress() bool {
	return printProgressBars && term.IsTerminal(int(os.Stdout.Fd()))
}

// wrappedRepo wraps oras.Target to show a progress bar on Push() operations.
type wrappedRepo struct {
	oras.Target
	progress *mpb.Progress
}

func (w *wrappedRepo) Push(ctx context.Context, expected ocispec.Descriptor, content io.Reader) error {
	shortDigest := expected.Digest.Encoded()[0:8]
	bar := w.progress.New(expected.Size,
		mpb.BarStyle().Lbound("|").Filler("=").Tip(">").Padding("-").Rbound("|"),
		mpb.PrependDecorators(
			decor.Name("Copying "+shortDigest),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.Counters(decor.SizeB1024(0), "% .1f / % .1f"), fmt.Sprintf("%-9s", FormatBytes(expected.Size))),
			decor.OnComplete(decor.Name(" | "), " | "),
			decor.OnComplete(decor.EwmaSpeed(decor.SizeB1024(0), "% .2f", 60), "done"),
		),
		mpb.BarFillerOnComplete("|"),
	)
	proxyReader := bar.ProxyReader(content)
	defer proxyReader.Close()

	return w.Target.Push(ctx, expected, proxyReader)
}

// WrapTarget wraps an oras.Target so that calls to Push print a progress bar.
// If output is configured to not print progress bars, this is a no-op.
func WrapTarget(wrap oras.Target) (oras.Target, *ProgressLogger) {
	if !shouldPrintProgress() {
		return wrap, &ProgressLogger{stdout}
	}
	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(180*time.Millisecond),
	)
	return &wrappedRepo{
		Target:   wrap,
		progress: p,
	}, &ProgressLogger{p}
}

func WrapReadCloser(size int64, rc io.ReadCloser) (io.ReadCloser, *ProgressLogger) {
	if !shouldPrintProgress() {
		return rc, &ProgressLogger{stdout}
	}

	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(180*time.Millisecond),
	)
	bar := p.New(size,
		mpb.BarStyle().Lbound("|").Filler("=").Tip(">").Padding("-").Rbound("|"),
		mpb.PrependDecorators(
			decor.Name("Unpacking"),
		),
		mpb.AppendDecorators(
			decor.Counters(decor.SizeB1024(0), "% .1f / % .1f"),
			decor.Name(" | "),
			decor.EwmaSpeed(decor.SizeB1024(0), "% .2f", 60),
		),
		mpb.BarRemoveOnComplete(),
	)

	return bar.ProxyReader(rc), &ProgressLogger{p}
}
