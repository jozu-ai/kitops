package pack

import (
	"context"
	"kitops/pkg/lib/constants"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackOptions_Complete(t *testing.T) {
	options := &packOptions{}
	ctx := context.WithValue(context.Background(), constants.ConfigKey{}, "/home/user/.kitops")
	args := []string{"arg1"}

	err := options.complete(ctx, args)

	assert.NoError(t, err)
	assert.Equal(t, args[0], options.contextDir)
	assert.Equal(t, filepath.Join(args[0], constants.DefaultKitFileName), options.modelFile)
}

func TestPackOptions_RunPack(t *testing.T) {
	t.Skip("Skipping test for now")
	options := &packOptions{
		modelFile:  "Kitfile",
		contextDir: "/path/to/context",
	}

	err := runPack(context.Background(), options)

	assert.NoError(t, err)
}
