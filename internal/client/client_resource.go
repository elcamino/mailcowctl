package client

import (
	"context"
	"fmt"
)

type Resource struct {
	ID               int    `json:"id"`
	Description      string `json:"description,omitempty"`
	Domain           string `json:"domain,omitempty"`
	Kind             string `json:"kind,omitempty"`
	MultipleBookings any    `json:"multiple_bookings,omitempty"`
	Active           any    `json:"active,omitempty"`
}

type ResourceCreate struct {
	Description      string `json:"description"`
	Domain           string `json:"domain"`
	Kind             string `json:"kind"`
	MultipleBookings int    `json:"multiple_bookings"`
	Active           string `json:"active,omitempty"`
}

func (c *Client) ListResources(ctx context.Context) ([]Resource, error) {
	var list []Resource
	err := c.getList(ctx, "/get/resource/all", &list)
	return list, err
}

func (c *Client) GetResource(ctx context.Context, id int) (Resource, error) {
	list, err := c.ListResources(ctx)
	if err != nil {
		return Resource{}, err
	}
	for _, r := range list {
		if r.ID == id {
			return r, nil
		}
	}
	return Resource{}, fmt.Errorf("resource %d not found", id)
}

func (c *Client) CreateResource(ctx context.Context, req ResourceCreate) error {
	return c.postAction(ctx, "/add/resource", req)
}

func (c *Client) DeleteResource(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/resource", []int{id})
}
