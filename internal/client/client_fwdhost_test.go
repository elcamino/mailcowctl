package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListAndGetFwdhostByHost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/get/fwdhost/all" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"host":"203.0.113.5","source":"203.0.113.5","keep_spam":"no"}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	h, err := c.GetFwdhost(context.Background(), "203.0.113.5")
	if err != nil {
		t.Fatal(err)
	}
	if h.Host != "203.0.113.5" || h.KeepSpam != "no" {
		t.Fatalf("fwdhost = %+v", h)
	}
}

func TestCreateAndDeleteFwdhost(t *testing.T) {
	var createBody, deleteBody json.RawMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/add/fwdhost":
			_ = json.NewDecoder(r.Body).Decode(&createBody)
		case "/api/v1/delete/fwdhost":
			_ = json.NewDecoder(r.Body).Decode(&deleteBody)
		default:
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"type":"success","msg":["ok"]}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	if err := c.CreateFwdhost(context.Background(), FwdhostCreate{Hostname: "203.0.113.5", FilterSpam: BoolString(true)}); err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteFwdhost(context.Background(), "203.0.113.5"); err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(createBody, &got); err != nil {
		t.Fatal(err)
	}
	if got["hostname"] != "203.0.113.5" || got["filter_spam"] != "1" {
		t.Fatalf("create body = %s", createBody)
	}
	if string(deleteBody) != `["203.0.113.5"]` {
		t.Fatalf("delete body = %s, want [\"203.0.113.5\"]", deleteBody)
	}
}
