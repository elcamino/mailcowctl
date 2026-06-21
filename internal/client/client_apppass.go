package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// AppPassword response field json tags (name, mailbox, and the *_access protocol flags) are
// based on mailcow's documented column names but have NOT been verified against a live populated
// response; reconcile these tags against real data during migration. Note: -o json/-o yaml output
// is unaffected because it echoes raw API keys.
type AppPassword struct {
	ID      int    `json:"id"`
	Name    string `json:"name,omitempty"`
	Mailbox string `json:"mailbox,omitempty"`
	Active  any    `json:"active,omitempty"`
	IMAP    any    `json:"imap_access,omitempty"`
	SMTP    any    `json:"smtp_access,omitempty"`
	POP3    any    `json:"pop3_access,omitempty"`
	Sieve   any    `json:"sieve_access,omitempty"`
	DAV     any    `json:"dav_access,omitempty"`
	EAS     any    `json:"eas_access,omitempty"`
}

func isEnabled(v any) bool {
	switch t := v.(type) {
	case string:
		return t == "1"
	case bool:
		return t
	default:
		return fmt.Sprint(v) == "1"
	}
}

func (a AppPassword) Protocols() string {
	var on []string
	for _, p := range []struct {
		name string
		val  any
	}{{"imap", a.IMAP}, {"smtp", a.SMTP}, {"pop3", a.POP3}, {"sieve", a.Sieve}, {"dav", a.DAV}, {"eas", a.EAS}} {
		if isEnabled(p.val) {
			on = append(on, p.name)
		}
	}
	return strings.Join(on, ",")
}

type AppPasswordCreate struct {
	Username   string   `json:"username"`
	AppName    string   `json:"app_name"`
	AppPasswd  string   `json:"app_passwd"`
	AppPasswd2 string   `json:"app_passwd2"`
	Protocols  []string `json:"protocols,omitempty"`
	Active     string   `json:"active,omitempty"`
}

func (c *Client) ListAppPasswords(ctx context.Context, mailbox string) ([]AppPassword, error) {
	var list []AppPassword
	err := c.getList(ctx, "/get/app-passwd/all/"+url.PathEscape(mailbox), &list)
	return list, err
}

func (c *Client) GetAppPassword(ctx context.Context, mailbox string, id int) (AppPassword, error) {
	list, err := c.ListAppPasswords(ctx, mailbox)
	if err != nil {
		return AppPassword{}, err
	}
	for _, a := range list {
		if a.ID == id {
			return a, nil
		}
	}
	return AppPassword{}, fmt.Errorf("app password %d not found for %s", id, mailbox)
}

func (c *Client) CreateAppPassword(ctx context.Context, req AppPasswordCreate) error {
	return c.postAction(ctx, "/add/app-passwd", req)
}

func (c *Client) EditAppPassword(ctx context.Context, id int, attr map[string]any) error {
	return c.postAction(ctx, "/edit/app-passwd", editRequest{Attr: attr, Items: []int{id}})
}

func (c *Client) DeleteAppPassword(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/app-passwd", []int{id})
}
