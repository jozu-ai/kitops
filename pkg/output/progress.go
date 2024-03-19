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
	"math"
	"time"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"oras.land/oras-go/v2"
)

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
func WrapTarget(wrap oras.Target) oras.Target {
	if !printProgressBars {
		return wrap
	}
	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(180*time.Millisecond),
	)
	return &wrappedRepo{
		Target:   wrap,
		progress: p,
	}
}

func WaitProgress(t oras.Target) {
	if wrapper, ok := t.(*wrappedRepo); ok {
		wrapper.progress.Wait()
	}
}

func FormatBytes(i int64) string {
	if i == 0 {
		return "0 B  "
	}

	if i < 1024 {
		// Catch bytes to avoid printing fractional amounts of bytes e.g. 123.0 bytes
		return fmt.Sprintf("%d B  ", i)
	}

	suffixes := []string{"KiB", "MiB", "GiB", "TiB"}
	unit := float64(1024)

	size := float64(i) / unit
	for _, suffix := range suffixes {
		if size < unit {
			// Round down to the nearest tenth of a unit to avoid e.g. 1MiB - 1B = 1024KiB
			niceSize := math.Floor(size*10) / 10
			return fmt.Sprintf("%.1f %s", niceSize, suffix)
		}
		size = size / unit
	}

	// Fall back to printing whatever's left as PiB
	return fmt.Sprintf("%.1f %s", size, "PiB")
}
