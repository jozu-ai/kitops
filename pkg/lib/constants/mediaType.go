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
	ZstdCompression = "zstd"
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
