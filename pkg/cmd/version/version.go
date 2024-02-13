package version

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

// Default build-time variable.
// These values are overridden via ldflags
var (
	Version   = "unknown"
	GitCommit = "unknown"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
)

func NewCmdVersion() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display the version information for jmm",
		Long: `The version command prints detailed version information for the jmm CLI tool,
including the current version of the tool, the Git commit that the version was built from, 
the build time, and the version of Go it was compiled with. This can be useful for debugging 
or verifying that you are running the expected version of jmm.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\nCommit: %s\nBuilt: %s\nGo version: %s\n", Version, GitCommit, BuildTime, GoVersion)
		},
	}
	return cmd
}

func init() {}
