package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListFiltersFiltersByMailbox(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/get/filters/all" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"id":1,"username":"a@x.org","filter_type":"prefilter"},{"id":2,"username":"b@x.org"}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	list, err := c.ListFilters(context.Background(), "a@x.org")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != 1 {
		t.Fatalf("filtered = %+v", list)
	}
}

func TestCreateFilterBody(t *testing.T) {
	var body json.RawMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&body)
		_, _ = w.Write([]byte(`[{"type":"success","msg":["ok"]}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	err := c.CreateFilter(context.Background(), FilterCreate{
		Username: "a@x.org", ScriptDesc: "vac", ScriptData: "require \"vacation\";", FilterType: "prefilter", Active: BoolString(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(body), `{"username":"a@x.org","script_desc":"vac","script_data":"require \"vacation\";","filter_type":"prefilter","active":"1"}`; got != want {
		t.Fatalf("body = %s, want %s", got, want)
	}
}
