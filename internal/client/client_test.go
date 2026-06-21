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

func TestDecodeActionTreatsDangerAsErrorDespiteHTTP200(t *testing.T) {
	err := DecodeActionResponse([]byte(`[{"type":"danger","msg":["domain already exists"]}]`))
	if err == nil {
		t.Fatal("DecodeActionResponse returned nil error, want APIError")
	}
	var apiErr *APIError
	if !AsAPIError(err, &apiErr) {
		t.Fatalf("error %T is not APIError", err)
	}
	if got := apiErr.Error(); got != "domain already exists" {
		t.Fatalf("error = %q, want mailcow message", got)
	}
}

func TestPostSendsXAPIKeyAndReturnsLogicalError(t *testing.T) {
	var sawKey string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawKey = r.Header.Get("X-API-Key")
		if r.URL.Path != "/api/v1/add/domain" {
			t.Fatalf("path = %s, want /api/v1/add/domain", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"type":"danger","msg":["read-only key"]}]`))
	}))
	defer server.Close()

	c, err := New(Options{Host: server.URL + "/", APIKey: "secret"})
	if err != nil {
		t.Fatal(err)
	}

	err = c.CreateDomain(context.Background(), DomainCreate{
		Domain: "example.org",
		Active: BoolString(true),
	})
	if err == nil {
		t.Fatal("CreateDomain returned nil error, want logical API error")
	}
	if sawKey != "secret" {
		t.Fatalf("X-API-Key = %q, want secret", sawKey)
	}
}

func TestDomainCreateAndEditBodies(t *testing.T) {
	var bodies [][]byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var raw json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			t.Fatal(err)
		}
		bodies = append(bodies, raw)
		_, _ = w.Write([]byte(`[{"type":"success","msg":["ok"]}]`))
	}))
	defer server.Close()

	c, err := New(Options{Host: server.URL, APIKey: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.CreateDomain(context.Background(), DomainCreate{
		Domain: "example.org",
		Quota:  10240,
		Active: BoolString(true),
	}); err != nil {
		t.Fatal(err)
	}
	if err := c.EditDomain(context.Background(), "example.org", map[string]any{
		"description": "new",
		"active":      BoolString(false),
	}); err != nil {
		t.Fatal(err)
	}

	if got, want := string(bodies[0]), `{"domain":"example.org","quota":10240,"active":"1"}`; got != want {
		t.Fatalf("create body = %s, want %s", got, want)
	}
	if got, want := string(bodies[1]), `{"attr":{"active":"0","description":"new"},"items":["example.org"]}`; got != want {
		t.Fatalf("edit body = %s, want %s", got, want)
	}
}

func TestAliasDeleteResolvesAddressToNumericID(t *testing.T) {
	var deleteBody json.RawMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/get/alias/all":
			_, _ = w.Write([]byte(`[{"id":123,"address":"alias@example.org","goto":"me@example.org","active":"1"}]`))
		case "/api/v1/delete/alias":
			if err := json.NewDecoder(r.Body).Decode(&deleteBody); err != nil {
				t.Fatal(err)
			}
			_, _ = w.Write([]byte(`[{"type":"success","msg":["deleted"]}]`))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c, err := New(Options{Host: server.URL, APIKey: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteAlias(context.Background(), "alias@example.org"); err != nil {
		t.Fatal(err)
	}

	if got, want := string(deleteBody), `[123]`; got != want {
		t.Fatalf("delete body = %s, want %s", got, want)
	}
}
