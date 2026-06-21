package client

import (
	"context"
	"fmt"
	"strings"
)

type Filter struct {
	ID         int    `json:"id"`
	Username   string `json:"username,omitempty"`
	ScriptDesc string `json:"script_desc,omitempty"`
	ScriptData string `json:"script_data,omitempty"`
	FilterType string `json:"filter_type,omitempty"`
	Active     any    `json:"active,omitempty"`
}

type FilterCreate struct {
	Username   string `json:"username"`
	ScriptDesc string `json:"script_desc"`
	ScriptData string `json:"script_data"`
	FilterType string `json:"filter_type"`
	Active     string `json:"active,omitempty"`
}

func (c *Client) ListFilters(ctx context.Context, mailbox string) ([]Filter, error) {
	var list []Filter
	if err := c.getList(ctx, "/get/filters/all", &list); err != nil {
		return nil, err
	}
	if mailbox == "" {
		return list, nil
	}
	filtered := list[:0]
	for _, f := range list {
		if strings.EqualFold(f.Username, mailbox) {
			filtered = append(filtered, f)
		}
	}
	return filtered, nil
}

func (c *Client) GetFilter(ctx context.Context, id int) (Filter, error) {
	list, err := c.ListFilters(ctx, "")
	if err != nil {
		return Filter{}, err
	}
	for _, f := range list {
		if f.ID == id {
			return f, nil
		}
	}
	return Filter{}, fmt.Errorf("filter %d not found", id)
}

func (c *Client) CreateFilter(ctx context.Context, req FilterCreate) error {
	return c.postAction(ctx, "/add/filter", req)
}

func (c *Client) EditFilter(ctx context.Context, id int, attr map[string]any) error {
	return c.postAction(ctx, "/edit/filter", editRequest{Attr: attr, Items: []int{id}})
}

func (c *Client) DeleteFilter(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/filter", []int{id})
}
