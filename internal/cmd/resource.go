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

func newResourceCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "resource", Short: "Manage calendar resources"}
	cmd.AddCommand(resourceListCmd(opts), resourceGetCmd(opts), resourceCreateCmd(opts), resourceDeleteCmd(opts))
	return cmd
}

func resourceListCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List calendar resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			list, err := c.ListResources(opts.context())
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, list)
		},
	}
}

func resourceGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a calendar resource",
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
			r, err := c.GetResource(opts.context(), id)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, r)
		},
	}
}

func resourceCreateCmd(opts *rootOptions) *cobra.Command {
	var description, domain, kind string
	var multipleBookings int
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a calendar resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			if description == "" || domain == "" || kind == "" {
				return fmt.Errorf("--description, --domain, and --kind are required")
			}
			req := client.ResourceCreate{
				Description: description, Domain: domain, Kind: kind, MultipleBookings: multipleBookings,
				Active: client.BoolString(true),
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
			return c.CreateResource(opts.context(), req)
		},
	}
	cmd.Flags().StringVar(&description, "description", "", "resource name")
	cmd.Flags().StringVar(&domain, "domain", "", "domain")
	cmd.Flags().StringVar(&kind, "kind", "", "location, group, or thing")
	cmd.Flags().IntVar(&multipleBookings, "multiple-bookings", 0, "multiple bookings value (0 disables)")
	cmd.Flags().Bool("active", false, "set active")
	cmd.Flags().Bool("inactive", false, "set inactive")
	return cmd
}

func resourceDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a calendar resource",
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
			return c.DeleteResource(opts.context(), id)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
