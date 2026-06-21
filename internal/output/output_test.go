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

package output

import "testing"

func TestMaskSecret(t *testing.T) {
	cases := []struct {
		short, full, want string
	}{
		{"abc...", "abcdef", "abc..."},
		{"", "abcdef", "***"},
		{"", "", ""},
	}
	for _, c := range cases {
		if got := maskSecret(c.short, c.full); got != c.want {
			t.Fatalf("maskSecret(%q,%q) = %q, want %q", c.short, c.full, got, c.want)
		}
	}
}
