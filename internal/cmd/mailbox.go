// SPDX-License-Identifier: GPL-3.0-or-later
//
// Copyright (C) 2026 Tobias von Dewitz
//
// This file is part of mailcowctl.
//
// mailcowctl is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// mailcowctl is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with mailcowctl. If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"
	"strconv"

	"github.com/elcamino/mailcowctl/internal/client"
	"github.com/elcamino/mailcowctl/internal/output"
	"github.com/spf13/cobra"
)

func newMailboxCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "mailbox", Short: "Manage mailboxes"}
	cmd.AddCommand(mailboxListCmd(opts), mailboxGetCmd(opts), mailboxCreateCmd(opts), mailboxEditCmd(opts), mailboxDeleteCmd(opts))
	return cmd
}

func mailboxListCmd(opts *rootOptions) *cobra.Command {
	var domain string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List mailboxes",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			mailboxes, err := c.ListMailboxes(opts.context(), domain)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, mailboxes)
		},
	}
	cmd.Flags().StringVar(&domain, "domain", "", "filter by domain")
	return cmd
}

func mailboxGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <address>",
		Short: "Get a mailbox",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			mailbox, err := c.GetMailbox(opts.context(), args[0])
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, mailbox)
		},
	}
}

func mailboxCreateCmd(opts *rootOptions) *cobra.Command {
	var name, passwordEnv string
	var quota int
	var passwordStdin, forcePasswordChange bool
	cmd := &cobra.Command{
		Use:   "create <address>",
		Short: "Create a mailbox",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			local, domain, err := splitMailbox(args[0])
			if err != nil {
				return err
			}
			password, err := passwordFromFlags(cmd, opts.in, passwordEnv, passwordStdin)
			if err != nil {
				return err
			}
			active, err := changedBool(cmd, "active", "inactive")
			if err != nil {
				return err
			}
			req := client.MailboxCreate{
				LocalPart:     local,
				Domain:        domain,
				Name:          name,
				Password:      password,
				Password2:     password,
				Active:        client.BoolString(true),
				ForcePwUpdate: client.BoolString(forcePasswordChange),
				TLSEnforceIn:  client.BoolString(false),
				TLSEnforceOut: client.BoolString(false),
			}
			if quota > 0 {
				req.Quota = strconv.Itoa(quota)
			}
			if active != nil {
				req.Active = client.BoolString(*active)
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.CreateMailbox(opts.context(), req)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "full name")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().IntVar(&quota, "quota", 0, "mailbox quota in MiB")
	cmd.Flags().BoolVar(&passwordStdin, "password-stdin", false, "read password from stdin")
	cmd.Flags().StringVar(&passwordEnv, "password-env", "", "read password from environment variable")
	cmd.Flags().Bool("active", false, "set mailbox active")
	cmd.Flags().Bool("inactive", false, "set mailbox inactive")
	cmd.Flags().BoolVar(&forcePasswordChange, "force-password-change", false, "force password change at next login")
	return cmd
}

func mailboxEditCmd(opts *rootOptions) *cobra.Command {
	var name string
	var quota int
	var passwordStdin bool
	cmd := &cobra.Command{
		Use:   "edit <address>",
		Short: "Edit a mailbox",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			attrs := map[string]any{}
			if cmd.Flags().Changed("name") {
				attrs["name"] = name
			}
			if cmd.Flags().Changed("quota") {
				attrs["quota"] = strconv.Itoa(quota)
			}
			active, err := changedBool(cmd, "active", "inactive")
			if err != nil {
				return err
			}
			if active != nil {
				attrs["active"] = client.BoolString(*active)
			}
			if passwordStdin {
				password, err := passwordFromFlags(cmd, opts.in, "", true)
				if err != nil {
					return err
				}
				attrs["password"] = password
				attrs["password2"] = password
			}
			if len(attrs) == 0 {
				return fmt.Errorf("no changes specified")
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.EditMailbox(opts.context(), args[0], attrs)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "full name")
	cmd.Flags().IntVar(&quota, "quota", 0, "mailbox quota in MiB")
	cmd.Flags().Bool("active", false, "set mailbox active")
	cmd.Flags().Bool("inactive", false, "set mailbox inactive")
	cmd.Flags().BoolVar(&passwordStdin, "password-stdin", false, "read password from stdin")
	return cmd
}

func mailboxDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <address>",
		Short: "Delete a mailbox",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := confirm(opts.in, opts.errOut, yes, args[0]); err != nil {
				return err
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.DeleteMailbox(opts.context(), args[0])
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
