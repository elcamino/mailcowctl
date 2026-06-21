package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tob/mailcowctl/internal/client"
	"github.com/tob/mailcowctl/internal/output"
)

func newDomainCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "domain", Short: "Manage domains"}
	cmd.AddCommand(domainListCmd(opts), domainGetCmd(opts), domainCreateCmd(opts), domainEditCmd(opts), domainDeleteCmd(opts))
	return cmd
}

func domainListCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List domains",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			domains, err := c.ListDomains(opts.context())
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, domains)
		},
	}
}

func domainGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <domain>",
		Short: "Get a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			domain, err := c.GetDomain(opts.context(), args[0])
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, domain)
		},
	}
}

func domainCreateCmd(opts *rootOptions) *cobra.Command {
	var req client.DomainCreate
	cmd := &cobra.Command{
		Use:   "create <domain>",
		Short: "Create a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			active, err := changedBool(cmd, "active", "inactive")
			if err != nil {
				return err
			}
			req.Domain = args[0]
			if active != nil {
				req.Active = client.BoolString(*active)
			} else {
				req.Active = client.BoolString(true)
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.CreateDomain(opts.context(), req)
		},
	}
	domainCreateEditFlags(cmd, &req)
	return cmd
}

func domainEditCmd(opts *rootOptions) *cobra.Command {
	var req client.DomainCreate
	cmd := &cobra.Command{
		Use:   "edit <domain>",
		Short: "Edit a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			attrs := map[string]any{}
			if cmd.Flags().Changed("description") {
				attrs["description"] = req.Description
			}
			if cmd.Flags().Changed("mailboxes") {
				attrs["mailboxes"] = req.Mailboxes
			}
			if cmd.Flags().Changed("aliases") {
				attrs["aliases"] = req.Aliases
			}
			if cmd.Flags().Changed("quota") {
				attrs["quota"] = req.Quota
			}
			if cmd.Flags().Changed("default-quota") {
				attrs["defquota"] = req.DefaultQuota
			}
			if cmd.Flags().Changed("max-quota") {
				attrs["maxquota"] = req.MaxQuota
			}
			active, err := changedBool(cmd, "active", "inactive")
			if err != nil {
				return err
			}
			if active != nil {
				attrs["active"] = client.BoolString(*active)
			}
			if len(attrs) == 0 {
				return fmt.Errorf("no changes specified")
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.EditDomain(opts.context(), args[0], attrs)
		},
	}
	domainCreateEditFlags(cmd, &req)
	return cmd
}

func domainDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <domain>",
		Short: "Delete a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := confirm(opts.in, opts.errOut, yes, args[0]); err != nil {
				return err
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.DeleteDomain(opts.context(), args[0])
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}

func domainCreateEditFlags(cmd *cobra.Command, req *client.DomainCreate) {
	cmd.Flags().StringVar(&req.Description, "description", "", "domain description")
	cmd.Flags().IntVar(&req.Mailboxes, "mailboxes", 0, "maximum mailboxes")
	cmd.Flags().IntVar(&req.Aliases, "aliases", 0, "maximum aliases")
	cmd.Flags().IntVar(&req.Quota, "quota", 0, "domain quota in MiB")
	cmd.Flags().IntVar(&req.DefaultQuota, "default-quota", 0, "default mailbox quota in MiB")
	cmd.Flags().IntVar(&req.MaxQuota, "max-quota", 0, "maximum mailbox quota in MiB")
	cmd.Flags().Bool("active", false, "set domain active")
	cmd.Flags().Bool("inactive", false, "set domain inactive")
}
