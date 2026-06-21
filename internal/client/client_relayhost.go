package client

import (
	"context"
	"fmt"
)

type Relayhost struct {
	ID              int    `json:"id"`
	Hostname        string `json:"hostname,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	PasswordShort   string `json:"password_short,omitempty"`
	Active          any    `json:"active,omitempty"`
	UsedByDomains   string `json:"used_by_domains,omitempty"`
	UsedByMailboxes string `json:"used_by_mailboxes,omitempty"`
}

type RelayhostCreate struct {
	Hostname string `json:"hostname"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Active   string `json:"active,omitempty"`
}

func (c *Client) ListRelayhosts(ctx context.Context) ([]Relayhost, error) {
	var list []Relayhost
	err := c.getList(ctx, "/get/relayhost/all", &list)
	return list, err
}

func (c *Client) GetRelayhost(ctx context.Context, id int) (Relayhost, error) {
	list, err := c.ListRelayhosts(ctx)
	if err != nil {
		return Relayhost{}, err
	}
	for _, h := range list {
		if h.ID == id {
			return h, nil
		}
	}
	return Relayhost{}, fmt.Errorf("relayhost %d not found", id)
}

func (c *Client) CreateRelayhost(ctx context.Context, req RelayhostCreate) error {
	return c.postAction(ctx, "/add/relayhost", req)
}

func (c *Client) DeleteRelayhost(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/relayhost", []int{id})
}
