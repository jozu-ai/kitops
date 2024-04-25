/*
Copyright © 2024 Jozu.com
*/
package login

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"kitops/pkg/cmd/options"
	"kitops/pkg/lib/constants"
	"kitops/pkg/output"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"oras.land/oras-go/v2/registry/remote/auth"
)

const (
	shortDesc = `Log in to an OCI registry`
	longDesc  = `Log in to a specified OCI-compatible registry. Credentials are saved and used
automatically for future CLI operations`

	example = `# Login to ghcr.io
kit login ghcr.io -u github_user -p personal_token`
)

type loginOptions struct {
	options.NetworkOptions
	registry          string
	configHome        string
	credential        auth.Credential
	username          string
	password          string
	passwordFromStdIn bool
}

func LoginCommand() *cobra.Command {
	opts := &loginOptions{}

	cmd := &cobra.Command{
		Use:     "login [flags] [REGISTRY]",
		Short:   shortDesc,
		Long:    longDesc,
		Example: example,
		RunE:    runCommand(opts),
	}

	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().StringVarP(&opts.username, "username", "u", "", "registry username")
	cmd.Flags().StringVarP(&opts.password, "password", "p", "", "registry password or token")
	cmd.Flags().BoolVar(&opts.passwordFromStdIn, "password-stdin", false, "read password from stdin")
	opts.AddNetworkFlags(cmd)

	return cmd
}

func runCommand(opts *loginOptions) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := opts.complete(cmd.Context(), args); err != nil {
			return output.Fatalf("Invalid arguments: %s", err)
		}

		err := login(cmd.Context(), opts)
		if err != nil {
			return output.Fatalln(err)
		}
		return nil
	}
}

func (opts *loginOptions) complete(ctx context.Context, args []string) error {
	opts.registry = args[0]
	configHome, ok := ctx.Value(constants.ConfigKey{}).(string)
	if !ok {
		return fmt.Errorf("default config path not set on command context")
	}
	opts.configHome = configHome

	if opts.password != "" && opts.username != "" {
		output.Infof("Warning: Using --password via CLI is insecure. Consider using --password-stdin instead")
	}

	username := opts.username
	password := opts.password
	if opts.passwordFromStdIn {
		if password != "" {
			return fmt.Errorf("cannot use both --password and --password-stdin")
		}
		readPass, err := readPasswordFromStdin()
		if err != nil {
			return err
		} else if readPass == "" {
			return fmt.Errorf("failed to read password from stdin: got empty string")
		}
		password = readPass
	}

	if password == "" {
		// Prompt for password (and username, if necessary)
		var err error
		if username == "" {
			username, err = promptForInput("Username: ", false)
			if err != nil {
				return err
			}
		}
		password, err = promptForInput("Password: ", true)
		if err != nil {
			return err
		}
		opts.credential = auth.Credential{
			Username: username,
			Password: password,
		}
	} else {
		// If username is empty, assume password is an OAuth token
		if username == "" {
			opts.credential = auth.Credential{
				RefreshToken: password,
			}
		} else {
			opts.credential = auth.Credential{
				Username: username,
				Password: password,
			}
		}
	}

	return nil
}

func readPasswordFromStdin() (string, error) {
	passwd, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read password from standard input")
	}
	return strings.TrimSpace(string(passwd)), err
}

func promptForInput(prompt string, isSensitive bool) (string, error) {
	var bytes []byte
	var err error
	if !term.IsTerminal(int(syscall.Stdin)) {
		return "", fmt.Errorf("attempting to read input from non-terminal")
	}

	fmt.Print(prompt)
	if isSensitive {
		bytes, err = term.ReadPassword(int(syscall.Stdin))
		fmt.Print("\n")
	} else {
		reader := bufio.NewReader(os.Stdin)
		bytes, err = reader.ReadBytes('\n')
	}
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return strings.TrimSpace(string(bytes)), nil
}
