package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tob/mailcowctl/internal/client"
	"github.com/tob/mailcowctl/internal/output"
)

func newTlspolicyCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "tlspolicy", Short: "Manage TLS policy maps"}
	cmd.AddCommand(tlspolicyListCmd(opts), tlspolicyGetCmd(opts), tlspolicyCreateCmd(opts), tlspolicyDeleteCmd(opts))
	return cmd
}

func tlspolicyListCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List TLS policy maps",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			list, err := c.ListTlsPolicies(opts.context())
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, list)
		},
	}
}

func tlspolicyGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a TLS policy map",
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
			p, err := c.GetTlsPolicy(opts.context(), id)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, p)
		},
	}
}

func tlspolicyCreateCmd(opts *rootOptions) *cobra.Command {
	var dest, policy, parameters string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a TLS policy map",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dest == "" || policy == "" {
				return fmt.Errorf("--dest and --policy are required")
			}
			req := client.TlsPolicyCreate{Dest: dest, Policy: policy, Parameters: parameters, Active: client.BoolString(true)}
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
			return c.CreateTlsPolicy(opts.context(), req)
		},
	}
	cmd.Flags().StringVar(&dest, "dest", "", "destination domain")
	cmd.Flags().StringVar(&policy, "policy", "", "none, may, encrypt, dane, dane-only, fingerprint, verify, or secure")
	cmd.Flags().StringVar(&parameters, "parameters", "", "optional policy parameters")
	cmd.Flags().Bool("active", false, "set active")
	cmd.Flags().Bool("inactive", false, "set inactive")
	return cmd
}

func tlspolicyDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a TLS policy map",
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
			return c.DeleteTlsPolicy(opts.context(), id)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
