package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListPolicyBothKinds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/get/policy_wl_domain/x.org":
			_, _ = w.Write([]byte(`{}`))
		case "/api/v1/get/policy_bl_domain/x.org":
			_, _ = w.Write([]byte(`[{"prefid":20,"object":"x.org","value":"*@spam.com"}]`))
		default:
			t.Fatalf("path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	items, err := c.ListPolicy(context.Background(), "x.org", "both")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Kind != "bl" || items[0].PrefID != 20 {
		t.Fatalf("items = %+v", items)
	}
}

func TestCreateAndDeletePolicyBodies(t *testing.T) {
	var bodies [][]byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var raw json.RawMessage
		_ = json.NewDecoder(r.Body).Decode(&raw)
		bodies = append(bodies, raw)
		_, _ = w.Write([]byte(`[{"type":"success","msg":["ok"]}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	if err := c.CreatePolicy(context.Background(), PolicyCreate{Domain: "x.org", ObjectList: "bl", ObjectFrom: "*@spam.com"}); err != nil {
		t.Fatal(err)
	}
	if err := c.DeletePolicy(context.Background(), 20); err != nil {
		t.Fatal(err)
	}
	if got, want := string(bodies[0]), `{"domain":"x.org","object_list":"bl","object_from":"*@spam.com"}`; got != want {
		t.Fatalf("create body = %s, want %s", got, want)
	}
	if got, want := string(bodies[1]), `[20]`; got != want {
		t.Fatalf("delete body = %s, want %s", got, want)
	}
}
