package build

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/assert"
)

func TestNewCmdBuild(t *testing.T) {
	cmd := NewCmdBuild()

	assert.NotNil(t, cmd)
	assert.Equal(t, "build", cmd.Use)
	assert.Equal(t, shortDesc, cmd.Short)
	assert.Equal(t, longDesc, cmd.Long)
}

func TestBuildOptions_Complete(t *testing.T) {
	options := &BuildOptions{}
	cmd := &cobra.Command{}
	args := []string{"arg1"}

	err := options.Complete(cmd, args)

	assert.NoError(t, err)
	assert.Equal(t, args[0], options.ContextDir)
	assert.Equal(t, options.ContextDir+"/"+DEFAULT_MODEL_FILE, options.ModelFile)
}

func TestBuildOptions_Validate(t *testing.T) {
	options := &BuildOptions{
		ModelFile:  "Jozufile",
		ContextDir: "/path/to/context",
	}

	err := options.Validate()

	assert.NoError(t, err)
}

func TestBuildOptions_RunBuild(t *testing.T) {
	t.Skip("Skipping test for now")
	options := &BuildOptions{
		ModelFile:  "Jozufile",
		ContextDir: "/path/to/context",
	}

	err := options.RunBuild()

	assert.NoError(t, err)
}

func TestBuildFlags_ToOptions(t *testing.T) {
	flags := &BuildFlags{
		ModelFile: "Jozufile",
	}

	options, err := flags.ToOptions()

	assert.NoError(t, err)
	assert.Equal(t, flags.ModelFile, options.ModelFile)
}

func TestBuildFlags_AddFlags(t *testing.T) {
	flags := &BuildFlags{}
	cmd := &cobra.Command{}

	flags.AddFlags(cmd)

	assert.NotNil(t, cmd.Flags().Lookup("file"))
}

func TestNewBuildFlags(t *testing.T) {
	flags := NewBuildFlags()

	assert.NotNil(t, flags)
}
