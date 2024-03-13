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
	"fmt"
	"os"
	"strings"
)

var (
	printDebug = false
)

func SetDebug(debug bool) {
	printDebug = debug
}

func Infoln(s any) {
	fmt.Println(s)
}

func Infof(s string, args ...any) {
	// Avoid printing incomplete lines
	if !strings.HasSuffix(s, "\n") {
		s = s + "\n"
	}
	fmt.Printf(s, args...)
}

func Errorln(s any) {
	fmt.Fprintln(os.Stderr, s)
}

func Errorf(s string, args ...any) {
	// Avoid printing incomplete lines
	if !strings.HasSuffix(s, "\n") {
		s = s + "\n"
	}
	fmt.Fprintf(os.Stderr, s, args...)
}

func Fatalln(s any) {
	Errorln(s)
	os.Exit(1)
}

func Fatalf(s string, args ...any) {
	Errorf(s, args...)
	os.Exit(1)
}

func Debugln(s any) {
	if printDebug {
		fmt.Println(s)
	}
}

func Debugf(s string, args ...any) {
	if !printDebug {
		return
	}
	// Avoid printing incomplete lines
	if !strings.HasSuffix(s, "\n") {
		s = s + "\n"
	}
	fmt.Printf(s, args...)
}
