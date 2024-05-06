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
	"regexp"
)

const (
	ConfigType    = "config"
	ModelType     = "model"
	ModelPartType = "modelpart"
	DatasetType   = "dataset"
	CodeType      = "code"
)

const (
	NoneCompression = "none"
	GzipCompression = "gzip"
)

var mediaTypeRegexp = regexp.MustCompile(`^application/vnd.kitops.modelkit.(\w+).v1.tar(?:\+(\w+))?`)

type MediaType struct {
	BaseType    string
	Compression string
}

var ModelConfigMediaType = MediaType{
	BaseType: ConfigType,
}

func (t MediaType) String() string {
	if t.BaseType == ConfigType {
		return "application/vnd.kitops.modelkit.config.v1+json"
	}
	if t.Compression == NoneCompression {
		return fmt.Sprintf("application/vnd.kitops.modelkit.%s.v1.tar", t.BaseType)
	}
	return fmt.Sprintf("application/vnd.kitops.modelkit.%s.v1.tar+%s", t.BaseType, t.Compression)
}

func ParseMediaType(s string) MediaType {
	if s == "application/vnd.kitops.modelkit.config.v1+json" {
		return MediaType{
			BaseType: ConfigType,
		}
	}
	match := mediaTypeRegexp.FindStringSubmatch(s)
	if match == nil {
		return MediaType{}
	}
	mediaType := MediaType{
		BaseType:    match[1],
		Compression: match[2],
	}
	if mediaType.Compression == "" {
		mediaType.Compression = NoneCompression
	}
	return mediaType
}
