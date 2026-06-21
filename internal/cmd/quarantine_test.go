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
