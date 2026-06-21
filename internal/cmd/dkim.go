package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tob/mailcowctl/internal/client"
	"github.com/tob/mailcowctl/internal/output"
)

func newDkimCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "dkim", Short: "Manage DKIM keys"}
	cmd.AddCommand(dkimGetCmd(opts), dkimAddCmd(opts), dkimDuplicateCmd(opts), dkimDeleteCmd(opts))
	return cmd
}

func dkimGetCmd(opts *rootOptions) *cobra.Command {
	var dns bool
	cmd := &cobra.Command{
		Use:   "get <domain>",
		Short: "Get the DKIM record for a domain (private key is never returned by the API)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			d, err := c.GetDkim(opts.context(), args[0])
			if err != nil {
				return err
			}
			if dns {
				fmt.Fprintf(opts.out, "%s._domainkey.%s. IN TXT \"%s\"\n", d.Selector, d.Domain, d.DkimTxt)
				return nil
			}
			return output.Render(opts.out, opts.output, d)
		},
	}
	cmd.Flags().BoolVar(&dns, "dns", false, "print the DNS TXT record line")
	return cmd
}

func dkimAddCmd(opts *rootOptions) *cobra.Command {
	var selector string
	var keySize int
	cmd := &cobra.Command{
		Use:   "add <domain>",
		Short: "Generate a DKIM key for a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.CreateDkim(opts.context(), client.DkimCreate{Domains: args[0], Selector: selector, KeySize: keySize})
		},
	}
	cmd.Flags().StringVar(&selector, "selector", "dkim", "DKIM selector")
	cmd.Flags().IntVar(&keySize, "key-size", 2048, "key size in bits")
	return cmd
}

func dkimDuplicateCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "duplicate <from-domain> <to-domain>",
		Short: "Copy a DKIM key to another domain on the same server",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.DuplicateDkim(opts.context(), client.DkimDuplicate{FromDomain: args[0], ToDomain: args[1]})
		},
	}
}

func dkimDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <domain>",
		Short: "Delete the DKIM key for a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := confirm(opts.in, opts.errOut, yes, args[0]); err != nil {
				return err
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.DeleteDkim(opts.context(), args[0])
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
