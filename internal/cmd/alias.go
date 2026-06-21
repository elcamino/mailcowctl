package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tob/mailcowctl/internal/client"
	"github.com/tob/mailcowctl/internal/output"
)

func newAliasCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "alias", Short: "Manage aliases"}
	cmd.AddCommand(aliasListCmd(opts), aliasGetCmd(opts), aliasCreateCmd(opts), aliasEditCmd(opts), aliasDeleteCmd(opts))
	return cmd
}

func aliasListCmd(opts *rootOptions) *cobra.Command {
	var domain string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			aliases, err := c.ListAliases(opts.context(), domain)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, aliases)
		},
	}
	cmd.Flags().StringVar(&domain, "domain", "", "filter by domain")
	return cmd
}

func aliasGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <address|id>",
		Short: "Get an alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			alias, err := c.GetAlias(opts.context(), args[0])
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, alias)
		},
	}
}

func aliasCreateCmd(opts *rootOptions) *cobra.Command {
	var gotoValue string
	var discard, toSpam, toHam, sogoHidden bool
	cmd := &cobra.Command{
		Use:   "create <address>",
		Short: "Create an alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := aliasTarget(gotoValue, discard, toSpam, toHam)
			if err != nil {
				return err
			}
			active, err := changedBool(cmd, "active", "inactive")
			if err != nil {
				return err
			}
			req := client.AliasCreate{
				Address:     args[0],
				Goto:        target,
				Active:      client.BoolString(true),
				SogoVisible: client.BoolString(!sogoHidden),
			}
			if active != nil {
				req.Active = client.BoolString(*active)
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.CreateAlias(opts.context(), req)
		},
	}
	aliasTargetFlags(cmd, &gotoValue, &discard, &toSpam, &toHam)
	cmd.Flags().Bool("active", false, "set alias active")
	cmd.Flags().Bool("inactive", false, "set alias inactive")
	cmd.Flags().BoolVar(&sogoHidden, "sogo-hidden", false, "hide alias from SOGo")
	return cmd
}

func aliasEditCmd(opts *rootOptions) *cobra.Command {
	var gotoValue string
	var discard, toSpam, toHam bool
	cmd := &cobra.Command{
		Use:   "edit <address|id>",
		Short: "Edit an alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			attrs := map[string]any{}
			if cmd.Flags().Changed("goto") || discard || toSpam || toHam {
				target, err := aliasTarget(gotoValue, discard, toSpam, toHam)
				if err != nil {
					return err
				}
				attrs["goto"] = target
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
			return c.EditAlias(opts.context(), args[0], attrs)
		},
	}
	aliasTargetFlags(cmd, &gotoValue, &discard, &toSpam, &toHam)
	cmd.Flags().Bool("active", false, "set alias active")
	cmd.Flags().Bool("inactive", false, "set alias inactive")
	return cmd
}

func aliasDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <address|id>",
		Short: "Delete an alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := confirm(opts.in, opts.errOut, yes, args[0]); err != nil {
				return err
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.DeleteAlias(opts.context(), args[0])
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}

func aliasTargetFlags(cmd *cobra.Command, gotoValue *string, discard, toSpam, toHam *bool) {
	cmd.Flags().StringVar(gotoValue, "goto", "", "comma-separated alias target addresses")
	cmd.Flags().BoolVar(discard, "discard", false, "silently discard mail")
	cmd.Flags().BoolVar(toSpam, "to-spam", false, "train mail as spam")
	cmd.Flags().BoolVar(toHam, "to-ham", false, "train mail as ham")
}

func aliasTarget(gotoValue string, discard, toSpam, toHam bool) (string, error) {
	var targets []string
	if gotoValue != "" {
		targets = append(targets, strings.TrimSpace(gotoValue))
	}
	if discard {
		targets = append(targets, "goto_null")
	}
	if toSpam {
		targets = append(targets, "goto_spam")
	}
	if toHam {
		targets = append(targets, "goto_ham")
	}
	if len(targets) != 1 {
		return "", errorsForAliasTarget(len(targets))
	}
	return targets[0], nil
}

func errorsForAliasTarget(count int) error {
	if count == 0 {
		return fmt.Errorf("one of --goto, --discard, --to-spam, or --to-ham is required")
	}
	return fmt.Errorf("use only one of --goto, --discard, --to-spam, or --to-ham")
}
