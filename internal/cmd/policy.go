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
