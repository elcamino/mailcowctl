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

	"github.com/elcamino/mailcowctl/internal/output"
	"github.com/spf13/cobra"
)

func newQuarantineCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "quarantine", Short: "Review and act on quarantined mail"}
	cmd.AddCommand(
		quarantineListCmd(opts),
		quarantineGetCmd(opts),
		quarantineReleaseCmd(opts),
		quarantineLearnHamCmd(opts),
		quarantineDeleteCmd(opts),
	)
	return cmd
}

func quarantineListCmd(opts *rootOptions) *cobra.Command {
	var rcpt string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List quarantined items",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			list, err := c.ListQuarantine(opts.context(), rcpt)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, list)
		},
	}
	cmd.Flags().StringVar(&rcpt, "rcpt", "", "filter by recipient address")
	return cmd
}

func quarantineGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a quarantined item",
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
			q, err := c.GetQuarantine(opts.context(), id)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, q)
		},
	}
}

func quarantineReleaseCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "release <id>",
		Short: "Release (deliver) a quarantined item",
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
			return c.ReleaseQuarantine(opts.context(), id)
		},
	}
}

func quarantineLearnHamCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "learn-ham <id>",
		Short: "Release a quarantined item and train it as ham",
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
			return c.LearnHamQuarantine(opts.context(), id)
		},
	}
}

func quarantineDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a quarantined item",
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
			return c.DeleteQuarantine(opts.context(), id)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
