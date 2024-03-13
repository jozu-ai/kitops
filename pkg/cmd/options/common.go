package options

import (
	"github.com/spf13/cobra"
)

// NetworkOptions represent common networking-related flags that are used by multiple commands.
// The flags should be added to the command via AddNetworkFlags before running.
type NetworkOptions struct {
	PlainHTTP bool
	TlsVerify bool
}

func (o *NetworkOptions) AddNetworkFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.PlainHTTP, "plain-http", false, "Use plain HTTP when connecting to remote registries")
	cmd.Flags().BoolVar(&o.TlsVerify, "tls-verify", true, "Require TLS and verify certificates when connecting to remote registries")
}
