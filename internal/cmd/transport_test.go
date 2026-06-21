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
	"io"
	"strings"
	"testing"
)

func TestTransportCreatePasswordSourceMutualExclusion(t *testing.T) {
	isolateConfig(t)
	cmd := NewRootCmd(strings.NewReader(""), io.Discard, io.Discard)
	cmd.SetArgs([]string{
		"--host", "http://127.0.0.1:1", "--api-key", "k",
		"transport", "create",
		"--destination", "example.org",
		"--nexthop", "[relay.example.org]:587",
		"--password-env", "SOMEVAR",
		"--password-stdin",
	})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when both --password-env and --password-stdin are set, got nil")
	}
	if !strings.Contains(err.Error(), "--password-env") || !strings.Contains(err.Error(), "--password-stdin") {
		t.Fatalf("expected error to mention both --password-env and --password-stdin, got: %s", err.Error())
	}
}
