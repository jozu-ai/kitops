package tag

import (
	"context"
	"fmt"
	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/errdef"
)

func RunTag(ctx context.Context, options *tagOptions) error {
	storageHome := constants.StoragePath(options.configHome)
	sourceStore, err := repo.NewLocalStore(storageHome, options.sourceRef)
	if err != nil {
		return fmt.Errorf("failed to open local storage: %w", err)
	}
	descriptor, err := oras.Resolve(ctx, sourceStore, options.sourceRef.Reference, oras.ResolveOptions{})
	if err != nil {
		if err == errdef.ErrNotFound {
			return fmt.Errorf("model %s not found", options.sourceRef.String())
		}
		return fmt.Errorf("error resolving model: %s", err)
	}
	if options.sourceRef.Registry == options.targetRef.Registry && options.sourceRef.Repository == options.targetRef.Repository {
		err = sourceStore.Tag(ctx, descriptor, options.targetRef.Reference)
		if err != nil {
			return fmt.Errorf("failed to tag reference %s: %w", options.targetRef, err)
		}
		return nil
	}
	// model kit is on a different registry and/or repo, copy the model to the target store
	targetStore, err := repo.NewLocalStore(storageHome, options.targetRef)
	if err != nil {
		return fmt.Errorf("failed to open local storage: %w", err)
	}
	_, err = oras.Copy(ctx, sourceStore, options.sourceRef.Reference, targetStore, options.targetRef.Reference, oras.CopyOptions{})
	if err != nil {
		return fmt.Errorf("failed to tag model: %w", err)
	}
	return nil
}
