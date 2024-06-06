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
	"io"
	"os"
)

type LogLevel int

const (
	LogLevelTrace LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func (l LogLevel) shouldPrint(atLevel LogLevel) bool {
	return l <= atLevel
}

func (l LogLevel) getOutput() io.Writer {
	switch l {
	case LogLevelError, LogLevelWarn:
		return stderr
	default:
		return stdout
	}
}

func (l LogLevel) getPrefix() string {
	if logLevel == LogLevelInfo {
		// By default, don't include log level in message
		return ""
	}
	switch l {
	case LogLevelTrace:
		return "[TRACE] "
	case LogLevelDebug:
		return "[DEBUG] "
	case LogLevelInfo:
		return "[INFO]  "
	case LogLevelWarn:
		return "[WARN]  "
	case LogLevelError:
		return "[ERROR] "
	}
	return ""
}

var (
	logLevel                  = LogLevelInfo
	progressStyle             = "plain"
	progressEnabled           = true
	stdout          io.Writer = os.Stdout
	stderr          io.Writer = os.Stderr
)

func SetLogLevel(level LogLevel) {
	logLevel = level
}

func SetLogLevelFromString(level string) error {
	switch level {
	case "trace":
		logLevel = LogLevelTrace
	case "debug":
		logLevel = LogLevelDebug
	case "info":
		logLevel = LogLevelInfo
	case "warn":
		logLevel = LogLevelWarn
	case "error":
		logLevel = LogLevelError
	default:
		return fmt.Errorf("invalid log level '%s'. Options are 'trace', 'debug', 'info', 'warn', 'error'", level)
	}
	return nil
}

func SetProgressBars(style string) {
	progressStyle = style
	progressEnabled = shouldPrintProgress()
}

func ProgressEnabled() bool {
	return progressEnabled
}

func SetOut(w io.Writer) {
	stdout = w
}

func SetErr(w io.Writer) {
	stderr = w
}
