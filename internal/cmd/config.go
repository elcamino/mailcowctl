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

	"github.com/elcamino/mailcowctl/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "config", Short: "Manage mailcowctl config"}
	cmd.AddCommand(&cobra.Command{
		Use:   "set-profile <name>",
		Short: "Set the current profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return config.SetCurrentProfile(opts.configPath, args[0])
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "current",
		Short: "Print the current profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, err := config.CurrentProfile(opts.configPath)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(opts.out, profile)
			return err
		},
	})
	return cmd
}
