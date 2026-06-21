package client

import (
	"context"
	"net/url"
)

type PolicyItem struct {
	PrefID int    `json:"prefid"`
	Object string `json:"object,omitempty"`
	Value  string `json:"value,omitempty"`
	Kind   string `json:"-"`
}

type PolicyCreate struct {
	Domain     string `json:"domain"`
	ObjectList string `json:"object_list"`
	ObjectFrom string `json:"object_from"`
}

func (c *Client) ListPolicy(ctx context.Context, domain, kind string) ([]PolicyItem, error) {
	var out []PolicyItem
	if kind == "" || kind == "wl" || kind == "both" {
		var wl []PolicyItem
		if err := c.getList(ctx, "/get/policy_wl_domain/"+url.PathEscape(domain), &wl); err != nil {
			return nil, err
		}
		for i := range wl {
			wl[i].Kind = "wl"
		}
		out = append(out, wl...)
	}
	if kind == "" || kind == "bl" || kind == "both" {
		var bl []PolicyItem
		if err := c.getList(ctx, "/get/policy_bl_domain/"+url.PathEscape(domain), &bl); err != nil {
			return nil, err
		}
		for i := range bl {
			bl[i].Kind = "bl"
		}
		out = append(out, bl...)
	}
	return out, nil
}

func (c *Client) CreatePolicy(ctx context.Context, req PolicyCreate) error {
	return c.postAction(ctx, "/add/domain-policy", req)
}

func (c *Client) DeletePolicy(ctx context.Context, prefid int) error {
	return c.postAction(ctx, "/delete/domain-policy", []int{prefid})
}
