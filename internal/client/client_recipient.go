package client

import (
	"context"
	"fmt"
)

type RecipientMap struct {
	ID     int    `json:"id"`
	Old    string `json:"recipient_map_old,omitempty"`
	New    string `json:"recipient_map_new,omitempty"`
	Active any    `json:"active,omitempty"`
}

type RecipientMapCreate struct {
	Old    string `json:"recipient_map_old"`
	New    string `json:"recipient_map_new"`
	Active string `json:"active,omitempty"`
}

func (c *Client) ListRecipientMaps(ctx context.Context) ([]RecipientMap, error) {
	var list []RecipientMap
	err := c.getList(ctx, "/get/recipient_map/all", &list)
	return list, err
}

func (c *Client) GetRecipientMap(ctx context.Context, id int) (RecipientMap, error) {
	list, err := c.ListRecipientMaps(ctx)
	if err != nil {
		return RecipientMap{}, err
	}
	for _, m := range list {
		if m.ID == id {
			return m, nil
		}
	}
	return RecipientMap{}, fmt.Errorf("recipient map %d not found", id)
}

func (c *Client) CreateRecipientMap(ctx context.Context, req RecipientMapCreate) error {
	return c.postAction(ctx, "/add/recipient_map", req)
}

func (c *Client) DeleteRecipientMap(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/recipient_map", []int{id})
}
