package tag

import (
	"context"
	"errors"
	"fmt"

	"kitops/pkg/lib/constants"
	"kitops/pkg/lib/repo/local"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/errdef"
)

// RunTag tags a local model or copies it to a new repository if the target is different.
func RunTag(ctx context.Context, options *tagOptions) error {
	storageHome := constants.StoragePath(options.configHome)

	sourceRepo, err := local.NewLocalRepo(storageHome, options.sourceRef)
	if err != nil {
		return fmt.Errorf("failed to open source local repository: %w", err)
	}

	descriptor, err := resolveSourceModel(ctx, sourceRepo, options.sourceRef)
	if err != nil {
		return err // The error is already wrapped inside resolveSourceModel
	}

	// If the source and target repository are the same, just update the tag.
	if isSameRepository(options.sourceRef, options.targetRef) {
		return tagInSameRepo(ctx, sourceRepo, descriptor, options.targetRef)
	}

	// If the target is a different repository, copy the manifest to the target.
	return tagInDifferentRepo(ctx, storageHome, sourceRepo, options.sourceRef, options.targetRef)
}

// resolveSourceModel resolves the source model descriptor.
func resolveSourceModel(ctx context.Context, repo local.Repository, sourceRef *Reference) (oras.Descriptor, error) {
	descriptor, err := oras.Resolve(ctx, repo, sourceRef.Reference, oras.ResolveOptions{})
	if err != nil {
		if errors.Is(err, errdef.ErrNotFound) {
			return oras.Descriptor{}, fmt.Errorf("model %s not found", sourceRef.String())
		}
		return oras.Descriptor{}, fmt.Errorf("error resolving model: %w", err)
	}
	return descriptor, nil
}

// isSameRepository checks if the source and target references point to the same repository.
func isSameRepository(sourceRef, targetRef *Reference) bool {
	return sourceRef.Registry == targetRef.Registry && sourceRef.Repository == targetRef.Repository
}

// tagInSameRepo tags the model in the same repository.
func tagInSameRepo(ctx context.Context, repo local.Repository, descriptor oras.Descriptor, targetRef *Reference) error {
	err := repo.Tag(ctx, descriptor, targetRef.Reference)
	if err != nil {
		return fmt.Errorf("failed to tag reference %s: %w", targetRef.String(), err)
	}
	return nil
}

// tagInDifferentRepo copies the manifest to a new repository and tags it.
func tagInDifferentRepo(ctx context.Context, storageHome string, sourceRepo local.Repository, sourceRef, targetRef *Reference) error {
	targetRepo, err := local.NewLocalRepo(storageHome, targetRef)
	if err != nil {
		return fmt.Errorf("failed to open target local repository: %w", err)
	}

	_, err = oras.Copy(ctx, sourceRepo, sourceRef.Reference, targetRepo, targetRef.Reference, oras.CopyOptions{})
	if err != nil {
		return fmt.Errorf("failed to tag model in different repository: %w", err)
	}
	return nil
}
