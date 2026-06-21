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

func TestListAppPasswordsEmptyObject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/get/app-passwd/all/a@x.org" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	list, err := c.ListAppPasswords(context.Background(), "a@x.org")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Fatalf("len = %d, want 0", len(list))
	}
}

func TestCreateAppPasswordBody(t *testing.T) {
	var body json.RawMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&body)
		_, _ = w.Write([]byte(`[{"type":"success","msg":["ok"]}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	err := c.CreateAppPassword(context.Background(), AppPasswordCreate{
		Username: "a@x.org", AppName: "phone", AppPasswd: "pw", AppPasswd2: "pw",
		Protocols: []string{"imap", "smtp"}, Active: BoolString(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(body), `{"username":"a@x.org","app_name":"phone","app_passwd":"pw","app_passwd2":"pw","protocols":["imap","smtp"],"active":"1"}`; got != want {
		t.Fatalf("body = %s, want %s", got, want)
	}
}

func TestAppPasswordProtocols(t *testing.T) {
	ap := AppPassword{IMAP: "1", SMTP: "1", POP3: "0"}
	if got := ap.Protocols(); got != "imap,smtp" {
		t.Fatalf("Protocols() = %q, want imap,smtp", got)
	}
}
