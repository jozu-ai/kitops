package remove

import (
	"context"
	"fmt"
	"kitops/pkg/lib/repo"
	"kitops/pkg/output"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/errdef"
	"oras.land/oras-go/v2/registry"
)

func removeModel(ctx context.Context, store repo.LocalStorage, ref *registry.Reference, forceDelete bool) (ocispec.Descriptor, error) {
	desc, err := oras.Resolve(ctx, store, ref.Reference, oras.ResolveOptions{})
	if err != nil {
		if err == errdef.ErrNotFound {
			return ocispec.DescriptorEmptyJSON, fmt.Errorf("model %s not found", ref.String())
		}
		return ocispec.DescriptorEmptyJSON, fmt.Errorf("error resolving model: %s", err)
	}

	// If reference passed in is a digest, remove the manifest ignoring any tags the manifest might have
	if err := ref.ValidateReferenceAsDigest(); err == nil || forceDelete {
		output.Debugf("Deleting manifest with digest %s", ref.Reference)
		if err := store.Delete(ctx, desc); err != nil {
			return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to delete model: %ws", err)
		}
		return desc, nil
	}

	tags, err := repo.GetTagsForDescriptor(ctx, store, desc)
	if err != nil {
		return ocispec.DescriptorEmptyJSON, err
	}
	if len(tags) <= 1 {
		output.Debugf("Deleting manifest tagged %s", ref.Reference)
		if err := store.Delete(ctx, desc); err != nil {
			return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to delete model: %w", err)
		}
	} else {
		output.Debugf("Found other tags for manifest: [%s]", strings.Join(tags, ", "))
		output.Debugf("Untagging %s", ref.Reference)
		if err := store.Untag(ctx, ref.Reference); err != nil {
			return ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to untag model: %w", err)
		}
	}

	return desc, nil
}
