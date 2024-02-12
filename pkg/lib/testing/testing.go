package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jmm/pkg/artifact"
	"jmm/pkg/lib/storage"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const ValidSize = int64(100)
const InvalidSize = int64(-1)

var TestingNotFoundError = errors.New("artifact not found")
var TestingInvalidSizeError = errors.New("invalid size")

type TestStore struct {
	// Repository for this store
	Repo string
	// Map of digest to Manifest, to simulate retrieval from e.g. disk
	Manifests map[digest.Digest]ocispec.Manifest
	// Map of digest to Config, to simulate retrieval from e.g. disk
	Configs map[digest.Digest]artifact.JozuFile
	// Index for the store
	Index *ocispec.Index
}

var _ storage.Store = (*TestStore)(nil)

// Fetch simulates fetching a blob from the store, given a descriptor. If the object does
// not exist, returns TestingNotFoundError. To simulate mismatched size between the descriptor's
// 'size' field and the size of the object, set the descriptor's size to InvalidSize
func (s *TestStore) Fetch(_ context.Context, desc ocispec.Descriptor) ([]byte, error) {
	for digest, manifest := range s.Manifests {
		if digest == desc.Digest {
			if desc.Size == InvalidSize {
				return nil, TestingInvalidSizeError
			}
			jsonBytes, err := json.Marshal(manifest)
			if err != nil {
				return nil, fmt.Errorf("testing -- unexpected error: failed to marshal manifest: %w", err)
			}
			return jsonBytes, nil
		}
	}
	for digest, config := range s.Configs {
		if digest == desc.Digest {
			if desc.Size == InvalidSize {
				return nil, TestingInvalidSizeError
			}
			jsonBytes, err := json.Marshal(config)
			if err != nil {
				return nil, fmt.Errorf("testing -- unexpected error: failed to marshal config: %w", err)
			}
			return jsonBytes, nil
		}
	}
	return nil, TestingNotFoundError
}

// ParseIndexJson simulates reading the index.json for the store. If an index json does not
// exist, returns TestingNotFoundError
func (s *TestStore) ParseIndexJson() (*ocispec.Index, error) {
	if s.Index != nil {
		return s.Index, nil
	}
	return nil, TestingNotFoundError
}

func (*TestStore) TagModel(ocispec.Descriptor, string) error {
	return fmt.Errorf("tag model is not implemented for testing")
}

// SaveModel is not yet implemented!
func (*TestStore) SaveModel(*artifact.Model, string) (*ocispec.Descriptor, error) {
	return nil, fmt.Errorf("save model is not implemented for testing")
}

func (t *TestStore) GetRepository() string {
	return t.Repo
}
