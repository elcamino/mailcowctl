package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListSyncJobsFiltersByMailbox(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/get/syncjobs/all/no_log" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"id":2,"user2":"a@x.org","host1":"h"},{"id":3,"user2":"b@x.org","host1":"h"}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	jobs, err := c.ListSyncJobs(context.Background(), "b@x.org")
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 || jobs[0].ID != 3 {
		t.Fatalf("filtered = %+v", jobs)
	}
}

func TestCreateSyncJobBody(t *testing.T) {
	var body json.RawMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&body)
		_, _ = w.Write([]byte(`[{"type":"success","msg":["ok"]}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	err := c.CreateSyncJob(context.Background(), SyncJobCreate{
		Username: "a@x.org", Host1: "mail.old", Port1: 993, User1: "a@old", Password1: "pw",
		Enc1: "SSL", MinsInterval: 20, Active: BoolString(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}
	if got["username"] != "a@x.org" || got["host1"] != "mail.old" || got["active"] != "1" {
		t.Fatalf("body = %s", body)
	}
}

func TestEditSyncJobBody(t *testing.T) {
	var body json.RawMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&body)
		_, _ = w.Write([]byte(`[{"type":"success","msg":["ok"]}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	if err := c.EditSyncJob(context.Background(), 7, map[string]any{"active": "0"}); err != nil {
		t.Fatal(err)
	}
	if got, want := string(body), `{"attr":{"active":"0"},"items":[7]}`; got != want {
		t.Fatalf("edit body = %s, want %s", got, want)
	}
}
