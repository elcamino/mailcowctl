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

func newRelayhostCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "relayhost", Short: "Manage sender-dependent relay hosts"}
	cmd.AddCommand(relayhostListCmd(opts), relayhostGetCmd(opts), relayhostCreateCmd(opts), relayhostDeleteCmd(opts))
	return cmd
}

func relayhostListCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List relay hosts (passwords masked; use -o json for full secrets)",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			list, err := c.ListRelayhosts(opts.context())
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, list)
		},
	}
}

func relayhostGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a relay host",
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
			h, err := c.GetRelayhost(opts.context(), id)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, h)
		},
	}
}

func relayhostCreateCmd(opts *rootOptions) *cobra.Command {
	var hostname, username, passwordEnv string
	var passwordStdin bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a relay host",
		RunE: func(cmd *cobra.Command, args []string) error {
			if hostname == "" {
				return fmt.Errorf("--hostname is required")
			}
			password := ""
			if passwordEnv != "" || passwordStdin {
				var err error
				password, err = passwordFromFlags(cmd, opts.in, passwordEnv, passwordStdin)
				if err != nil {
					return err
				}
			}
			req := client.RelayhostCreate{Hostname: hostname, Username: username, Password: password, Active: client.BoolString(true)}
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
			return c.CreateRelayhost(opts.context(), req)
		},
	}
	cmd.Flags().StringVar(&hostname, "hostname", "", "relay host, e.g. smtp.example.org:587")
	cmd.Flags().StringVar(&username, "username", "", "auth username (optional)")
	cmd.Flags().StringVar(&passwordEnv, "password-env", "", "env var holding the auth password")
	cmd.Flags().BoolVar(&passwordStdin, "password-stdin", false, "read the auth password from stdin")
	cmd.Flags().Bool("active", false, "set active")
	cmd.Flags().Bool("inactive", false, "set inactive")
	return cmd
}

func relayhostDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a relay host",
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
			return c.DeleteRelayhost(opts.context(), id)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
