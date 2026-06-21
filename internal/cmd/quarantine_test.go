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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestQuarantineDeleteRequiresConfirmation(t *testing.T) {
	isolateConfig(t)
	cmd := NewRootCmd(strings.NewReader("no\n"), io.Discard, io.Discard)
	cmd.SetArgs([]string{"quarantine", "delete", "5"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when confirmation is declined")
	}
}

func TestQuarantineReleasePostsActionBody(t *testing.T) {
	isolateConfig(t)
	var body json.RawMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/edit/qitem" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		_, _ = w.Write([]byte(`[{"type":"success","msg":["ok"]}]`))
	}))
	defer server.Close()

	cmd := NewRootCmd(strings.NewReader(""), io.Discard, io.Discard)
	cmd.SetArgs([]string{"--host", server.URL, "--api-key", "k", "quarantine", "release", "33"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if got, want := string(body), `{"attr":{"action":"release"},"items":[33]}`; got != want {
		t.Fatalf("release body = %s, want %s", got, want)
	}
}
