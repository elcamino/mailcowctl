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

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetListTreatsEmptyObjectAsEmptySlice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c, err := New(Options{Host: server.URL, APIKey: "k"})
	if err != nil {
		t.Fatal(err)
	}
	out := []map[string]any{{"keep": "me"}}
	if err := c.getList(context.Background(), "/get/filters/all", &out); err != nil {
		t.Fatalf("getList returned error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("len(out) = %d, want 0", len(out))
	}
}

func TestGetListDecodesArray(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":1},{"id":2}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	var out []map[string]any
	if err := c.getList(context.Background(), "/x", &out); err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("len(out) = %d, want 2", len(out))
	}
}
