package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDkimSetsDomainAndParsesFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/get/dkim/example.org" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"dkim_selector":"dkim","length":"2048","privkey":"","pubkey":"PUB","dkim_txt":"v=DKIM1;..."}`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	d, err := c.GetDkim(context.Background(), "example.org")
	if err != nil {
		t.Fatal(err)
	}
	if d.Domain != "example.org" || d.Selector != "dkim" || d.DkimTxt != "v=DKIM1;..." {
		t.Fatalf("unexpected dkim: %+v", d)
	}
}

func TestDkimCreateAndDuplicateBodies(t *testing.T) {
	var bodies [][]byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var raw json.RawMessage
		_ = json.NewDecoder(r.Body).Decode(&raw)
		bodies = append(bodies, raw)
		_, _ = w.Write([]byte(`[{"type":"success","msg":["ok"]}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	if err := c.CreateDkim(context.Background(), DkimCreate{Domains: "example.org", Selector: "dkim", KeySize: 2048}); err != nil {
		t.Fatal(err)
	}
	if err := c.DuplicateDkim(context.Background(), DkimDuplicate{FromDomain: "a.org", ToDomain: "b.org"}); err != nil {
		t.Fatal(err)
	}
	if got, want := string(bodies[0]), `{"domains":"example.org","dkim_selector":"dkim","key_size":2048}`; got != want {
		t.Fatalf("create body = %s, want %s", got, want)
	}
	if got, want := string(bodies[1]), `{"from_domain":"a.org","to_domain":"b.org"}`; got != want {
		t.Fatalf("duplicate body = %s, want %s", got, want)
	}
}
