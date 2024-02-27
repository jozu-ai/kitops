package repo

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"kitops/pkg/lib/constants"
	"os"
	"path"
	"path/filepath"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
)

type LocalStorage interface {
	GetRepo() string
	GetIndex() (*ocispec.Index, error)
	oras.Target
}

type LocalStore struct {
	storePath string
	repo      string
	*oci.Store
}

func GetAllLocalStores(storageRoot string) ([]LocalStorage, error) {
	subDirs, err := findStoragePaths(storageRoot)
	if err != nil {
		return nil, err
	}
	var stores []LocalStorage
	for _, subDir := range subDirs {
		// convert to forward slashes for repo
		repo := filepath.ToSlash(subDir)
		storePath := filepath.Join(storageRoot, subDir)
		ociStore, err := oci.New(storePath)
		if err != nil {
			return nil, err
		}
		localStore := &LocalStore{
			storePath: storePath,
			repo:      repo,
			Store:     ociStore,
		}
		stores = append(stores, localStore)
	}
	return stores, nil
}

func NewLocalStore(storageRoot string, ref *registry.Reference) (LocalStorage, error) {
	storePath := storageRoot
	repo := ""
	if ref != nil {
		repo = path.Join(ref.Registry, ref.Repository)
		storePath = filepath.Join(storePath, ref.Registry, ref.Repository)
	}
	store, err := oci.New(storePath)
	if err != nil {
		return nil, err
	}
	return &LocalStore{
		storePath: storePath,
		repo:      repo,
		Store:     store,
	}, nil
}

func (s *LocalStore) GetIndex() (*ocispec.Index, error) {
	return parseIndexJson(s.storePath)
}

func (s *LocalStore) GetRepo() string {
	return s.repo
}

func findStoragePaths(storageRoot string) ([]string, error) {
	var indexPaths []string
	err := filepath.WalkDir(storageRoot, func(file string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "index.json" && !info.IsDir() {
			dir := filepath.Dir(file)
			relDir, err := filepath.Rel(storageRoot, dir)
			if err != nil {
				return err
			}
			if relDir == "." {
				relDir = ""
			}
			indexPaths = append(indexPaths, relDir)
		}
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read local storage: %w", err)
	}
	return indexPaths, nil
}

func parseIndexJson(storageHome string) (*ocispec.Index, error) {
	indexBytes, err := os.ReadFile(constants.IndexJsonPath(storageHome))
	if err != nil {
		if os.IsNotExist(err) {
			return &ocispec.Index{}, nil
		}
		return nil, fmt.Errorf("failed to read index: %w", err)
	}

	index := &ocispec.Index{}
	if err := json.Unmarshal(indexBytes, index); err != nil {
		return nil, fmt.Errorf("failed to parse index: %w", err)
	}

	return index, nil
}
