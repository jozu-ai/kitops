package storage

import (
	"fmt"
	"path"
	"strings"

	"oras.land/oras-go/v2/registry"
)

// ParseReference parses a reference string into a Reference struct. If the
// reference does not include a registry (e.g. myrepo:mytag), the placeholder
// 'localhost' is used. Additional tags can be specified in a comma-separated
// list (e.g. myrepo:tag1,tag2,tag3)
func ParseReference(refString string) (ref *registry.Reference, extraTags []string, err error) {
	// References _must_ contain host; use localhost to mark local-only
	if !strings.Contains(refString, "/") {
		refString = fmt.Sprintf("localhost/%s", refString)
	}

	refAndTags := strings.Split(refString, ",")
	baseRef, err := registry.ParseReference(refAndTags[0])
	if err != nil {
		return nil, nil, err
	}
	return &baseRef, refAndTags[1:], nil
}

func StorageHome(configRoot string) string {
	return path.Join(configRoot, "storage")
}

func LocalStorePath(storageRoot string, ref *registry.Reference) string {
	return path.Join(storageRoot, ref.Registry, ref.Repository)
}
