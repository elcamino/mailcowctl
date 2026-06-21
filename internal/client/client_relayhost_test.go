package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListRelayhostsParsesPassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/get/relayhost/all" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"id":2,"hostname":"smtp.example.org:587","username":"u","password":"secret","password_short":"sec...","active":1}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	list, err := c.ListRelayhosts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Password != "secret" || list[0].PasswordShort != "sec..." {
		t.Fatalf("list = %+v", list)
	}
}

func TestCreateAndDeleteRelayhost(t *testing.T) {
	var createBody, deleteBody json.RawMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/add/relayhost":
			_ = json.NewDecoder(r.Body).Decode(&createBody)
		case "/api/v1/delete/relayhost":
			_ = json.NewDecoder(r.Body).Decode(&deleteBody)
		default:
			t.Fatalf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"type":"success","msg":["ok"]}]`))
	}))
	defer server.Close()

	c, _ := New(Options{Host: server.URL, APIKey: "k"})
	if err := c.CreateRelayhost(context.Background(), RelayhostCreate{
		Hostname: "smtp.example.org:587", Username: "u", Password: "p", Active: BoolString(true),
	}); err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteRelayhost(context.Background(), 2); err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(createBody, &got); err != nil {
		t.Fatal(err)
	}
	if got["hostname"] != "smtp.example.org:587" || got["username"] != "u" {
		t.Fatalf("create body = %s", createBody)
	}
	if string(deleteBody) != `[2]` {
		t.Fatalf("delete body = %s, want [2]", deleteBody)
	}
}
