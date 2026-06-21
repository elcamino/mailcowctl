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
