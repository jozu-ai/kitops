package build

import (
	"context"
	"kitops/pkg/lib/constants"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCmdBuild(t *testing.T) {
	cmd := BuildCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "build", cmd.Use)
	assert.Equal(t, shortDesc, cmd.Short)
	assert.Equal(t, longDesc, cmd.Long)
}

func TestBuildOptions_Complete(t *testing.T) {
	options := &buildOptions{}
	flags := &buildFlags{}
	ctx := context.WithValue(context.Background(), constants.ConfigKey{}, "/home/user/.kitops")
	args := []string{"arg1"}

	err := options.complete(ctx, flags, args)

	assert.NoError(t, err)
	assert.Equal(t, args[0], options.contextDir)
	assert.Equal(t, filepath.Join(args[0], constants.DefaultKitFileName), options.modelFile)
}

func TestBuildOptions_RunBuild(t *testing.T) {
	t.Skip("Skipping test for now")
	options := &buildOptions{
		modelFile:  "Kitfile",
		contextDir: "/path/to/context",
	}

	err := RunBuild(context.Background(), options)

	assert.NoError(t, err)
}
