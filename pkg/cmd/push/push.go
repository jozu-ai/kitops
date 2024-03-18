/*
Copyright © 2024 Jozu.com
*/
package push

import (
	"context"
	"fmt"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

func PushModel(ctx context.Context, localStore *oci.Store, remoteRegistry *remote.Registry, ref *registry.Reference) (ocispec.Descriptor, error) {
	repo, err := remoteRegistry.Repository(ctx, ref.Repository)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to read repository: %w", err)
	}

	desc, err := oras.Copy(ctx, localStore, ref.Reference, repo, ref.Reference, oras.DefaultCopyOptions)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to copy to remote: %w", err)
	}
	output.WaitProgress(trackedRepo)

	return desc, err
}
