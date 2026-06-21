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

func newFilterCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "filter", Short: "Manage sieve filters"}
	cmd.AddCommand(filterListCmd(opts), filterGetCmd(opts), filterCreateCmd(opts), filterEditCmd(opts), filterDeleteCmd(opts))
	return cmd
}

func filterListCmd(opts *rootOptions) *cobra.Command {
	var mailbox string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sieve filters",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			list, err := c.ListFilters(opts.context(), mailbox)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, list)
		},
	}
	cmd.Flags().StringVar(&mailbox, "mailbox", "", "filter by mailbox")
	return cmd
}

func filterGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a sieve filter",
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
			f, err := c.GetFilter(opts.context(), id)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, f)
		},
	}
}

func filterCreateCmd(opts *rootOptions) *cobra.Command {
	var mailbox, desc, filterType, scriptFile string
	var scriptStdin bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a sieve filter",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mailbox == "" || desc == "" {
				return fmt.Errorf("--mailbox and --desc are required")
			}
			if filterType != "prefilter" && filterType != "postfilter" {
				return fmt.Errorf("--type must be prefilter or postfilter")
			}
			script, err := scriptFromFlags(opts.in, scriptFile, scriptStdin)
			if err != nil {
				return err
			}
			req := client.FilterCreate{
				Username: mailbox, ScriptDesc: desc, ScriptData: script, FilterType: filterType,
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
			return c.CreateFilter(opts.context(), req)
		},
	}
	cmd.Flags().StringVar(&mailbox, "mailbox", "", "mailbox address")
	cmd.Flags().StringVar(&desc, "desc", "", "filter description")
	cmd.Flags().StringVar(&filterType, "type", "prefilter", "prefilter or postfilter")
	cmd.Flags().StringVar(&scriptFile, "script-file", "", "path to a sieve script file")
	cmd.Flags().BoolVar(&scriptStdin, "script-stdin", false, "read the sieve script from stdin")
	cmd.Flags().Bool("active", false, "set active")
	cmd.Flags().Bool("inactive", false, "set inactive")
	return cmd
}

func filterEditCmd(opts *rootOptions) *cobra.Command {
	var desc, filterType, scriptFile string
	var scriptStdin bool
	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit a sieve filter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("id must be numeric")
			}
			attrs := map[string]any{}
			stringFlagAttr(cmd, "desc", "script_desc", attrs)
			stringFlagAttr(cmd, "type", "filter_type", attrs)
			if scriptFile != "" || scriptStdin {
				script, err := scriptFromFlags(opts.in, scriptFile, scriptStdin)
				if err != nil {
					return err
				}
				attrs["script_data"] = script
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
			return c.EditFilter(opts.context(), id, attrs)
		},
	}
	cmd.Flags().StringVar(&desc, "desc", "", "filter description")
	cmd.Flags().StringVar(&filterType, "type", "", "prefilter or postfilter")
	cmd.Flags().StringVar(&scriptFile, "script-file", "", "path to a sieve script file")
	cmd.Flags().BoolVar(&scriptStdin, "script-stdin", false, "read the sieve script from stdin")
	cmd.Flags().Bool("active", false, "set active")
	cmd.Flags().Bool("inactive", false, "set inactive")
	return cmd
}

func filterDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a sieve filter",
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
			return c.DeleteFilter(opts.context(), id)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
