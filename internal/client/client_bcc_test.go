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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAndDeleteBcc(t *testing.T) {
	var createBody, deleteBody json.RawMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/add/bcc":
			_ = json.NewDecoder(r.Body).Decode(&createBody)
		case "/api/v1/delete/bcc":
			_ = json.NewDecoder(r.Body).Decode(&deleteBody)
		default:
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"type":"success","msg":["ok"]}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	if err := c.CreateBcc(context.Background(), BccCreate{
		LocalDest: "example.org", BccDest: "archive@example.org", Type: "rcpt", Active: BoolString(true),
	}); err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteBcc(context.Background(), 3); err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(createBody, &got); err != nil {
		t.Fatal(err)
	}
	if got["local_dest"] != "example.org" || got["bcc_dest"] != "archive@example.org" || got["type"] != "rcpt" {
		t.Fatalf("create body = %s", createBody)
	}
	if string(deleteBody) != `[3]` {
		t.Fatalf("delete body = %s, want [3]", deleteBody)
	}
}

func TestListBccsEmptyObject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()
	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	list, err := c.ListBccs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Fatalf("len = %d, want 0", len(list))
	}
}
