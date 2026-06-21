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
