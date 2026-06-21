package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/tob/mailcowctl/internal/client"
	"github.com/tob/mailcowctl/internal/config"
)

const (
	ExitAPIError    = 1
	ExitUsageError  = 2
	ExitConfigError = 3
)

type rootOptions struct {
	host       string
	apiKey     string
	profile    string
	output     string
	insecure   bool
	quiet      bool
	verbose    bool
	configPath string
	in         io.Reader
	out        io.Writer
	errOut     io.Writer
}

func Execute() {
	root := NewRootCmd(os.Stdin, os.Stdout, os.Stderr)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ExitCode(err))
	}
}

func NewRootCmd(in io.Reader, out, errOut io.Writer) *cobra.Command {
	opts := &rootOptions{in: in, out: out, errOut: errOut, output: "table"}
	root := &cobra.Command{
		Use:           "mailcowctl",
		Short:         "Manage mailcow domains, mailboxes, and aliases",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().StringVar(&opts.host, "host", "", "mailcow base URL")
	root.PersistentFlags().StringVar(&opts.apiKey, "api-key", "", "mailcow API key")
	root.PersistentFlags().StringVar(&opts.profile, "profile", "", "config profile")
	root.PersistentFlags().StringVarP(&opts.output, "output", "o", "table", "output format: table, json, yaml")
	root.PersistentFlags().BoolVar(&opts.insecure, "insecure", false, "skip TLS certificate verification")
	root.PersistentFlags().BoolVarP(&opts.quiet, "quiet", "q", false, "suppress non-essential messages")
	root.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "print extra diagnostics to stderr")
	root.PersistentFlags().StringVar(&opts.configPath, "config", "", "config file path")

	root.AddCommand(newVersionCmd(opts))
	root.AddCommand(newDomainCmd(opts))
	root.AddCommand(newMailboxCmd(opts))
	root.AddCommand(newAliasCmd(opts))
	root.AddCommand(newDkimCmd(opts))
	root.AddCommand(newSyncjobCmd(opts))
	root.AddCommand(newApppassCmd(opts))
	root.AddCommand(newFilterCmd(opts))
	root.AddCommand(newPolicyCmd(opts))
	root.AddCommand(newTransportCmd(opts))
	root.AddCommand(newConfigCmd(opts))
	root.AddCommand(newCompletionCmd(root))
	return root
}

func (o *rootOptions) client() (*client.Client, error) {
	cfg, err := config.Resolve(config.Inputs{
		Host:       o.host,
		APIKey:     o.apiKey,
		Profile:    o.profile,
		ConfigPath: o.configPath,
	})
	if err != nil {
		return nil, err
	}
	return client.New(client.Options{Host: cfg.Host, APIKey: cfg.APIKey, Insecure: o.insecure})
}

func (o *rootOptions) context() context.Context {
	return context.Background()
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	if config.IsAuthConfigError(err) {
		return ExitConfigError
	}
	var authErr *client.AuthError
	if errors.As(err, &authErr) {
		return ExitConfigError
	}
	var apiErr *client.APIError
	if errors.As(err, &apiErr) {
		return ExitAPIError
	}
	return ExitUsageError
}
