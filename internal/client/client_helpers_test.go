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

package client

import "testing"

func TestFindByID(t *testing.T) {
	type item struct {
		ID   int
		Name string
	}
	items := []item{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "two"},
	}

	got, err := findByID(items, 2, func(item item) int {
		return item.ID
	}, "thing")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "two" {
		t.Fatalf("item.Name = %q, want two", got.Name)
	}
}

func TestFindByIDNotFound(t *testing.T) {
	type item struct {
		ID int
	}
	items := []item{{ID: 1}}

	got, err := findByID(items, 3, func(item item) int {
		return item.ID
	}, "thing")
	if err == nil {
		t.Fatal("findByID returned nil error")
	}
	if err.Error() != "thing 3 not found" {
		t.Fatalf("err = %q, want thing 3 not found", err)
	}
	if got.ID != 0 {
		t.Fatalf("item.ID = %d, want zero value", got.ID)
	}
}
