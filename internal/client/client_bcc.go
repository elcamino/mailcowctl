package client

import (
	"context"
	"fmt"
)

type Bcc struct {
	ID        int    `json:"id"`
	Type      string `json:"type,omitempty"`
	LocalDest string `json:"local_dest,omitempty"`
	BccDest   string `json:"bcc_dest,omitempty"`
	Domain    string `json:"domain,omitempty"`
	Active    any    `json:"active,omitempty"`
}

type BccCreate struct {
	LocalDest string `json:"local_dest"`
	BccDest   string `json:"bcc_dest"`
	Type      string `json:"type"`
	Active    string `json:"active,omitempty"`
}

func (c *Client) ListBccs(ctx context.Context) ([]Bcc, error) {
	var list []Bcc
	err := c.getList(ctx, "/get/bcc/all", &list)
	return list, err
}

func (c *Client) GetBcc(ctx context.Context, id int) (Bcc, error) {
	list, err := c.ListBccs(ctx)
	if err != nil {
		return Bcc{}, err
	}
	for _, b := range list {
		if b.ID == id {
			return b, nil
		}
	}
	return Bcc{}, fmt.Errorf("bcc %d not found", id)
}

func (c *Client) CreateBcc(ctx context.Context, req BccCreate) error {
	return c.postAction(ctx, "/add/bcc", req)
}

func (c *Client) DeleteBcc(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/bcc", []int{id})
}
