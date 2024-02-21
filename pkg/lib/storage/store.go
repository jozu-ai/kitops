package storage

import (
	"context"
	"kitops/pkg/artifact"

	_ "crypto/sha256"
	_ "crypto/sha512"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type Store interface {
	SaveModel(model *artifact.Model, tag string) (*ocispec.Descriptor, error)
	TagModel(manifestDesc ocispec.Descriptor, tag string) error
	GetRepository() string
	ParseIndexJson() (*ocispec.Index, error)
	Fetch(context.Context, ocispec.Descriptor) ([]byte, error)
}
