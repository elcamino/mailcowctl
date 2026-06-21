package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tob/mailcowctl/internal/client"
	"github.com/tob/mailcowctl/internal/output"
)

func newBccCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "bcc", Short: "Manage BCC maps"}
	cmd.AddCommand(bccListCmd(opts), bccGetCmd(opts), bccCreateCmd(opts), bccDeleteCmd(opts))
	return cmd
}

func bccListCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List BCC maps",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			list, err := c.ListBccs(opts.context())
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, list)
		},
	}
}

func bccGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a BCC map",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("id must be numeric")
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			b, err := c.GetBcc(opts.context(), id)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, b)
		},
	}
}

func bccCreateCmd(opts *rootOptions) *cobra.Command {
	var localDest, bccDest, bccType string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a BCC map",
		RunE: func(cmd *cobra.Command, args []string) error {
			if localDest == "" || bccDest == "" {
				return fmt.Errorf("--local-dest and --bcc-dest are required")
			}
			if bccType != "sender" && bccType != "rcpt" {
				return fmt.Errorf("--type must be sender or rcpt")
			}
			req := client.BccCreate{
				LocalDest: localDest, BccDest: bccDest, Type: bccType,
				Active: client.BoolString(true),
			}
			active, err := changedBool(cmd, "active", "inactive")
			if err != nil {
				return err
			}
			if active != nil {
				req.Active = client.BoolString(*active)
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.CreateBcc(opts.context(), req)
		},
	}
	cmd.Flags().StringVar(&localDest, "local-dest", "", "local destination (domain or address)")
	cmd.Flags().StringVar(&bccDest, "bcc-dest", "", "BCC destination address")
	cmd.Flags().StringVar(&bccType, "type", "", "sender or rcpt")
	cmd.Flags().Bool("active", false, "set active")
	cmd.Flags().Bool("inactive", false, "set inactive")
	return cmd
}

func bccDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a BCC map",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("id must be numeric")
			}
			if err := confirm(opts.in, opts.errOut, yes, args[0]); err != nil {
				return err
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.DeleteBcc(opts.context(), id)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
