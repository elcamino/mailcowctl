package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tob/mailcowctl/internal/client"
	"github.com/tob/mailcowctl/internal/output"
)

func newApppassCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "apppass", Short: "Manage app passwords"}
	cmd.AddCommand(apppassListCmd(opts), apppassGetCmd(opts), apppassCreateCmd(opts), apppassEditCmd(opts), apppassDeleteCmd(opts))
	return cmd
}

func apppassListCmd(opts *rootOptions) *cobra.Command {
	var mailbox string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List app passwords for a mailbox (secrets are never returned)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mailbox == "" {
				return fmt.Errorf("--mailbox is required")
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			list, err := c.ListAppPasswords(opts.context(), mailbox)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, list)
		},
	}
	cmd.Flags().StringVar(&mailbox, "mailbox", "", "mailbox address")
	return cmd
}

func apppassGetCmd(opts *rootOptions) *cobra.Command {
	var mailbox string
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get an app password",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if mailbox == "" {
				return fmt.Errorf("--mailbox is required")
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("id must be numeric")
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			ap, err := c.GetAppPassword(opts.context(), mailbox, id)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, ap)
		},
	}
	cmd.Flags().StringVar(&mailbox, "mailbox", "", "mailbox address")
	return cmd
}

func apppassCreateCmd(opts *rootOptions) *cobra.Command {
	var mailbox, name, protocols, passwordEnv string
	var passwordStdin bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an app password",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mailbox == "" || name == "" {
				return fmt.Errorf("--mailbox and --name are required")
			}
			password, err := passwordFromFlags(cmd, opts.in, passwordEnv, passwordStdin)
			if err != nil {
				return err
			}
			req := client.AppPasswordCreate{
				Username: mailbox, AppName: name, AppPasswd: password, AppPasswd2: password,
				Protocols: splitCSV(protocols), Active: client.BoolString(true),
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
			return c.CreateAppPassword(opts.context(), req)
		},
	}
	cmd.Flags().StringVar(&mailbox, "mailbox", "", "mailbox address")
	cmd.Flags().StringVar(&name, "name", "", "app password name")
	cmd.Flags().StringVar(&protocols, "protocols", "imap,smtp", "comma-separated protocols")
	cmd.Flags().StringVar(&passwordEnv, "password-env", "", "env var holding the app password")
	cmd.Flags().BoolVar(&passwordStdin, "password-stdin", false, "read app password from stdin")
	cmd.Flags().Bool("active", false, "set active")
	cmd.Flags().Bool("inactive", false, "set inactive")
	return cmd
}

func apppassEditCmd(opts *rootOptions) *cobra.Command {
	var name, protocols string
	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit an app password",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("id must be numeric")
			}
			attrs := map[string]any{}
			if cmd.Flags().Changed("name") {
				attrs["app_name"] = name
			}
			if cmd.Flags().Changed("protocols") {
				attrs["protocols"] = splitCSV(protocols)
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
			return c.EditAppPassword(opts.context(), id, attrs)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "app password name")
	cmd.Flags().StringVar(&protocols, "protocols", "", "comma-separated protocols")
	cmd.Flags().Bool("active", false, "set active")
	cmd.Flags().Bool("inactive", false, "set inactive")
	return cmd
}

func apppassDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an app password",
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
			return c.DeleteAppPassword(opts.context(), id)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}

func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
