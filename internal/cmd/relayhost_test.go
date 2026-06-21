package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func relayhostTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/get/relayhost/all" {
			payload := `[{"id":1,"hostname":"relay.example.org:587","username":"user","password":"SUPERSECRET","password_short":"SUP...","active":1}]`
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, payload)
			return
		}
		http.NotFound(w, r)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestRelayhostListMasksPasswordInTable(t *testing.T) {
	isolateConfig(t)
	srv := relayhostTestServer(t)

	var out bytes.Buffer
	cmd := NewRootCmd(strings.NewReader(""), &out, io.Discard)
	cmd.SetArgs([]string{"--host", srv.URL, "--api-key", "k", "relayhost", "list"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	if strings.Contains(got, "SUPERSECRET") {
		t.Errorf("table output must not contain the full password SUPERSECRET, got:\n%s", got)
	}
	if !strings.Contains(got, "SUP...") {
		t.Errorf("table output must contain the short form SUP..., got:\n%s", got)
	}
}

func TestRelayhostListJSONIncludesFullPassword(t *testing.T) {
	isolateConfig(t)
	srv := relayhostTestServer(t)

	var out bytes.Buffer
	cmd := NewRootCmd(strings.NewReader(""), &out, io.Discard)
	cmd.SetArgs([]string{"--host", srv.URL, "--api-key", "k", "-o", "json", "relayhost", "list"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw := out.Bytes()
	var list []map[string]any
	if err := json.Unmarshal(raw, &list); err != nil {
		t.Fatalf("JSON output is not valid: %v\nraw: %s", err, raw)
	}
	if len(list) == 0 {
		t.Fatal("expected at least one relayhost in JSON output")
	}
	pw, _ := list[0]["password"].(string)
	if pw != "SUPERSECRET" {
		t.Errorf("JSON output must include the full password SUPERSECRET, got %q", pw)
	}
}
