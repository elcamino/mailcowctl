package client

import (
	"context"
	"fmt"
)

type Transport struct {
	ID            int    `json:"id"`
	Destination   string `json:"destination,omitempty"`
	Nexthop       string `json:"nexthop,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	PasswordShort string `json:"password_short,omitempty"`
	Active        any    `json:"active,omitempty"`
}

type TransportCreate struct {
	Destination string `json:"destination"`
	Nexthop     string `json:"nexthop"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Active      string `json:"active,omitempty"`
}

func (c *Client) ListTransports(ctx context.Context) ([]Transport, error) {
	var list []Transport
	err := c.getList(ctx, "/get/transport/all", &list)
	return list, err
}

func (c *Client) GetTransport(ctx context.Context, id int) (Transport, error) {
	list, err := c.ListTransports(ctx)
	if err != nil {
		return Transport{}, err
	}
	for _, t := range list {
		if t.ID == id {
			return t, nil
		}
	}
	return Transport{}, fmt.Errorf("transport %d not found", id)
}

func (c *Client) CreateTransport(ctx context.Context, req TransportCreate) error {
	return c.postAction(ctx, "/add/transport", req)
}

func (c *Client) DeleteTransport(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/transport", []int{id})
}
