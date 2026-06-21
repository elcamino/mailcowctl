package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tob/mailcowctl/internal/client"
	"github.com/tob/mailcowctl/internal/output"
)

func newFwdhostCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "fwdhost", Short: "Manage forwarding hosts"}
	cmd.AddCommand(fwdhostListCmd(opts), fwdhostGetCmd(opts), fwdhostCreateCmd(opts), fwdhostDeleteCmd(opts))
	return cmd
}

func fwdhostListCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List forwarding hosts",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			list, err := c.ListFwdhosts(opts.context())
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, list)
		},
	}
}

func fwdhostGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <host>",
		Short: "Get a forwarding host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			h, err := c.GetFwdhost(opts.context(), args[0])
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, h)
		},
	}
}

func fwdhostCreateCmd(opts *rootOptions) *cobra.Command {
	var hostname string
	var filterSpam bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a forwarding host",
		RunE: func(cmd *cobra.Command, args []string) error {
			if hostname == "" {
				return fmt.Errorf("--hostname is required")
			}
			req := client.FwdhostCreate{Hostname: hostname, FilterSpam: client.BoolString(filterSpam)}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.CreateFwdhost(opts.context(), req)
		},
	}
	cmd.Flags().StringVar(&hostname, "hostname", "", "forwarding host IP or hostname")
	cmd.Flags().BoolVar(&filterSpam, "filter-spam", false, "keep spam filtering for mail from this host")
	return cmd
}

func fwdhostDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <host>",
		Short: "Delete a forwarding host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := confirm(opts.in, opts.errOut, yes, args[0]); err != nil {
				return err
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.DeleteFwdhost(opts.context(), args[0])
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
