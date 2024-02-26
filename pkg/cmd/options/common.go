package options

import (
	"github.com/spf13/cobra"
)

type NetworkOptions struct {
	PlainHTTP bool
	TlsVerify bool
}

func (o *NetworkOptions) AddNetworkFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.PlainHTTP, "plain-http", false, "Use plain HTTP when connecting to remote registries")
	cmd.Flags().BoolVar(&o.TlsVerify, "tls-verify", true, "Require TLS and verify certificates when connecting to remote registries")
}
