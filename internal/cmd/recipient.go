package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tob/mailcowctl/internal/client"
	"github.com/tob/mailcowctl/internal/output"
)

func newRecipientCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "recipient", Short: "Manage recipient maps"}
	cmd.AddCommand(recipientListCmd(opts), recipientGetCmd(opts), recipientCreateCmd(opts), recipientDeleteCmd(opts))
	return cmd
}

func recipientListCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List recipient maps",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			list, err := c.ListRecipientMaps(opts.context())
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, list)
		},
	}
}

func recipientGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a recipient map",
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
			m, err := c.GetRecipientMap(opts.context(), id)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, m)
		},
	}
}

func recipientCreateCmd(opts *rootOptions) *cobra.Command {
	var oldAddr, newAddr string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a recipient map",
		RunE: func(cmd *cobra.Command, args []string) error {
			if oldAddr == "" || newAddr == "" {
				return fmt.Errorf("--old and --new are required")
			}
			req := client.RecipientMapCreate{Old: oldAddr, New: newAddr, Active: client.BoolString(true)}
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
			return c.CreateRecipientMap(opts.context(), req)
		},
	}
	cmd.Flags().StringVar(&oldAddr, "old", "", "original recipient")
	cmd.Flags().StringVar(&newAddr, "new", "", "rewritten recipient")
	cmd.Flags().Bool("active", false, "set active")
	cmd.Flags().Bool("inactive", false, "set inactive")
	return cmd
}

func recipientDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a recipient map",
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
			return c.DeleteRecipientMap(opts.context(), id)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
