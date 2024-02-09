package storage

import (
	"context"
	"jmm/pkg/artifact"

	_ "crypto/sha256"
	_ "crypto/sha512"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type Store interface {
	SaveModel(*artifact.Model) (*ocispec.Manifest, error)
	ParseIndexJson() (*ocispec.Index, error)
	Fetch(context.Context, ocispec.Descriptor) ([]byte, error)
}
