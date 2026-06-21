package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAndDeleteTlsPolicy(t *testing.T) {
	var createBody, deleteBody json.RawMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/add/tls-policy-map":
			_ = json.NewDecoder(r.Body).Decode(&createBody)
		case "/api/v1/delete/tls-policy-map":
			_ = json.NewDecoder(r.Body).Decode(&deleteBody)
		default:
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"type":"success","msg":["ok"]}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	if err := c.CreateTlsPolicy(context.Background(), TlsPolicyCreate{
		Dest: "example.org", Policy: "encrypt", Parameters: "", Active: BoolString(true),
	}); err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteTlsPolicy(context.Background(), 9); err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(createBody, &got); err != nil {
		t.Fatal(err)
	}
	if got["dest"] != "example.org" || got["policy"] != "encrypt" {
		t.Fatalf("create body = %s", createBody)
	}
	if string(deleteBody) != `[9]` {
		t.Fatalf("delete body = %s, want [9]", deleteBody)
	}
}
