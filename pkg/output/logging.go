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
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/vbauerster/mpb/v8"
)

func Infoln(s any) {
	fmt.Fprintln(stdout, s)
}

func Infof(s string, args ...any) {
	printFmt(stdout, s, args...)
}

func Errorln(s any) {
	fmt.Fprintln(stderr, s)
}

func Errorf(s string, args ...any) {
	printFmt(stderr, s, args...)
}

// Fatalln is the equivalent of Errorln except it returns a basic error to signal the command has failed
func Fatalln(s any) error {
	fmt.Fprintln(stderr, s)
	return errors.New("failed to run")
}

// Fatalf is the equivalent of Errorf except it returns a basic error to signal the command has failed
func Fatalf(s string, args ...any) error {
	printFmt(stderr, s, args...)
	return errors.New("failed to run")
}

func Debugln(s any) {
	if printDebug {
		fmt.Fprintln(stdout, s)
	}
}

func Debugf(s string, args ...any) {
	if printDebug {
		printFmt(stdout, s, args...)
	}
}

// ProgressLogger allows for printing info and debug lines while a progress bar
// is filling, and should be used instead of the standard output functions to prevent
// progress bars from removing log lines. Once the progress bar is done, the Wait()
// method should be called.
type ProgressLogger struct {
	output io.Writer
}

// Wait will call Wait() on the underlying mpb.Progress, if present. Otherwise,
// this is a no-op.
func (pw *ProgressLogger) Wait() {
	if progress, ok := pw.output.(*mpb.Progress); ok {
		progress.Wait()
	}
}

func (pw *ProgressLogger) Infoln(s any) {
	fmt.Fprintln(pw.output, s)
}

func (pw *ProgressLogger) Infof(s string, args ...any) {
	printFmt(pw.output, s, args...)
}

func (pw *ProgressLogger) Debugln(s any) {
	if printDebug {
		fmt.Fprintln(pw.output, s)
	}
}

func (pw *ProgressLogger) Debugf(s string, args ...any) {
	if printDebug {
		printFmt(pw.output, s, args...)
	}
}

func printFmt(w io.Writer, s string, args ...any) {
	// Avoid printing incomplete lines
	if !strings.HasSuffix(s, "\n") {
		s = s + "\n"
	}
	fmt.Fprintf(w, s, args...)
}
