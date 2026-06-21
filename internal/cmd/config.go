package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tob/mailcowctl/internal/config"
)

func newConfigCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "config", Short: "Manage mailcowctl config"}
	cmd.AddCommand(&cobra.Command{
		Use:   "set-profile <name>",
		Short: "Set the current profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return config.SetCurrentProfile(opts.configPath, args[0])
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "current",
		Short: "Print the current profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, err := config.CurrentProfile(opts.configPath)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(opts.out, profile)
			return err
		},
	})
	return cmd
}
