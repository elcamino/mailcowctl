package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListQuarantineFiltersByRcpt(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/get/quarantine/all" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"id":1,"rcpt":"a@x.org","sender":"s1","subject":"hi"},{"id":2,"rcpt":"b@x.org","sender":"s2","subject":"yo"}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	list, err := c.ListQuarantine(context.Background(), "b@x.org")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != 2 {
		t.Fatalf("filtered = %+v", list)
	}
}

func TestQuarantineActionAndDeleteBodies(t *testing.T) {
	var captured []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var raw json.RawMessage
		_ = json.NewDecoder(r.Body).Decode(&raw)
		captured = append(captured, r.URL.Path+" "+string(raw))
		_, _ = w.Write([]byte(`[{"type":"success","msg":["ok"]}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	if err := c.ReleaseQuarantine(context.Background(), 33); err != nil {
		t.Fatal(err)
	}
	if err := c.LearnHamQuarantine(context.Background(), 34); err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteQuarantine(context.Background(), 35); err != nil {
		t.Fatal(err)
	}

	want := []string{
		`/api/v1/edit/qitem {"attr":{"action":"release"},"items":[33]}`,
		`/api/v1/edit/qitem {"attr":{"action":"learnham"},"items":[34]}`,
		`/api/v1/delete/qitem {"items":[35]}`,
	}
	if len(captured) != len(want) {
		t.Fatalf("captured %d calls: %v", len(captured), captured)
	}
	for i := range want {
		if captured[i] != want[i] {
			t.Fatalf("call %d = %q, want %q", i, captured[i], want[i])
		}
	}
}
