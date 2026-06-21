package client

import (
	"context"
	"fmt"
)

type TlsPolicy struct {
	ID         int    `json:"id"`
	Dest       string `json:"dest,omitempty"`
	Policy     string `json:"policy,omitempty"`
	Parameters string `json:"parameters,omitempty"`
	Active     any    `json:"active,omitempty"`
}

type TlsPolicyCreate struct {
	Dest       string `json:"dest"`
	Policy     string `json:"policy"`
	Parameters string `json:"parameters,omitempty"`
	Active     string `json:"active,omitempty"`
}

func (c *Client) ListTlsPolicies(ctx context.Context) ([]TlsPolicy, error) {
	var list []TlsPolicy
	err := c.getList(ctx, "/get/tls-policy-map/all", &list)
	return list, err
}

func (c *Client) GetTlsPolicy(ctx context.Context, id int) (TlsPolicy, error) {
	list, err := c.ListTlsPolicies(ctx)
	if err != nil {
		return TlsPolicy{}, err
	}
	for _, p := range list {
		if p.ID == id {
			return p, nil
		}
	}
	return TlsPolicy{}, fmt.Errorf("tls policy %d not found", id)
}

func (c *Client) CreateTlsPolicy(ctx context.Context, req TlsPolicyCreate) error {
	return c.postAction(ctx, "/add/tls-policy-map", req)
}

func (c *Client) DeleteTlsPolicy(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/tls-policy-map", []int{id})
}
