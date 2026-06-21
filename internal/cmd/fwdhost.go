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

	"github.com/elcamino/mailcowctl/internal/client"
	"github.com/elcamino/mailcowctl/internal/output"
	"github.com/spf13/cobra"
)

func newFwdhostCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "fwdhost", Short: "Manage forwarding hosts"}
	cmd.AddCommand(fwdhostListCmd(opts), fwdhostGetCmd(opts), fwdhostCreateCmd(opts), fwdhostDeleteCmd(opts))
	return cmd
}

func fwdhostListCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List forwarding hosts",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			list, err := c.ListFwdhosts(opts.context())
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, list)
		},
	}
}

func fwdhostGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <host>",
		Short: "Get a forwarding host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			h, err := c.GetFwdhost(opts.context(), args[0])
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, h)
		},
	}
}

func fwdhostCreateCmd(opts *rootOptions) *cobra.Command {
	var hostname string
	var filterSpam bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a forwarding host",
		RunE: func(cmd *cobra.Command, args []string) error {
			if hostname == "" {
				return fmt.Errorf("--hostname is required")
			}
			req := client.FwdhostCreate{Hostname: hostname, FilterSpam: client.BoolString(filterSpam)}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.CreateFwdhost(opts.context(), req)
		},
	}
	cmd.Flags().StringVar(&hostname, "hostname", "", "forwarding host IP or hostname")
	cmd.Flags().BoolVar(&filterSpam, "filter-spam", false, "keep spam filtering for mail from this host")
	return cmd
}

func fwdhostDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <host>",
		Short: "Delete a forwarding host",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := confirm(opts.in, opts.errOut, yes, args[0]); err != nil {
				return err
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.DeleteFwdhost(opts.context(), args[0])
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
