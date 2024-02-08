package cmd

import (
	"os/user"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	t.Run("ToOptions", func(t *testing.T) {
		t.Run("WithConfigHome", func(t *testing.T) {
			f := &RootFlags{
				ConfigHome: "/path/to/config",
			}
			options, err := f.ToOptions()
			assert.NoError(t, err)
			assert.Equal(t, "/path/to/config", options.ConfigHome)
		})
	})
	t.Run("Complete", func(t *testing.T) {
		t.Run("WithConfigHome", func(t *testing.T) {
			o := &RootOptions{
				ConfigHome: "/path/to/config",
			}
			err := o.Complete()
			assert.NoError(t, err)
			assert.Equal(t, "/path/to/config", o.ConfigHome)
		})

		t.Run("WithoutConfigHome", func(t *testing.T) {
			currentUser, err := user.Current()
			assert.NoError(t, err)
			o := &RootOptions{}

			err = o.Complete()
			assert.NoError(t, err)
			assert.Equal(t, filepath.Join(currentUser.HomeDir, ".jozu"), o.ConfigHome)
			assert.Equal(t, filepath.Join(currentUser.HomeDir, ".jozu"), viper.GetString("config"))

		})
	})

}
