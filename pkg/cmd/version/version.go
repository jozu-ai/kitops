package version

import (
	"kitops/pkg/output"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	shortDesc = `Display the version information for the CLI`
	longDesc  = `The version command prints detailed version information.

This information includes the current version of the tool, the Git commit that 
the version was built from, the build time, and the version of Go it was 
compiled with.`
)

// Default build-time variable.
// These values are overridden via ldflags
var (
	Version   = "unknown"
	GitCommit = "unknown"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
)

func VersionCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "version",
		Short: shortDesc,
		Long:  longDesc,
		Run: func(cmd *cobra.Command, args []string) {
			output.Infof("Version: %s\nCommit: %s\nBuilt: %s\nGo version: %s\n", Version, GitCommit, BuildTime, GoVersion)
		},
	}
	return cmd
}

func init() {}
