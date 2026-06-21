package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tob/mailcowctl/internal/client"
	"github.com/tob/mailcowctl/internal/output"
)

func newPolicyCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "policy", Short: "Manage domain anti-spam allow/block lists"}
	cmd.AddCommand(policyListCmd(opts), policyAddCmd(opts), policyDeleteCmd(opts))
	return cmd
}

func policyListCmd(opts *rootOptions) *cobra.Command {
	var kind string
	cmd := &cobra.Command{
		Use:   "list <domain>",
		Short: "List allow/block list entries for a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch kind {
			case "wl", "bl", "both":
			default:
				return fmt.Errorf("--kind must be wl, bl, or both")
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			items, err := c.ListPolicy(opts.context(), args[0], kind)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, items)
		},
	}
	cmd.Flags().StringVar(&kind, "kind", "both", "wl, bl, or both")
	return cmd
}

func policyAddCmd(opts *rootOptions) *cobra.Command {
	var kind, value string
	cmd := &cobra.Command{
		Use:   "add <domain>",
		Short: "Add an allow (wl) or block (bl) list entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if kind != "wl" && kind != "bl" {
				return fmt.Errorf("--kind must be wl or bl")
			}
			if value == "" {
				return fmt.Errorf("--value is required")
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.CreatePolicy(opts.context(), client.PolicyCreate{Domain: args[0], ObjectList: kind, ObjectFrom: value})
		},
	}
	cmd.Flags().StringVar(&kind, "kind", "", "wl or bl")
	cmd.Flags().StringVar(&value, "value", "", "address or pattern, e.g. *@spam.com")
	return cmd
}

func policyDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <prefid>",
		Short: "Delete an allow/block list entry by prefid",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prefid, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("prefid must be numeric")
			}
			if err := confirm(opts.in, opts.errOut, yes, args[0]); err != nil {
				return err
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.DeletePolicy(opts.context(), prefid)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
