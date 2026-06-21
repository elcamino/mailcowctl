package client

import (
	"context"
	"fmt"
	"strings"
)

type Fwdhost struct {
	Host     string `json:"host,omitempty"`
	Source   string `json:"source,omitempty"`
	KeepSpam string `json:"keep_spam,omitempty"`
}

type FwdhostCreate struct {
	Hostname   string `json:"hostname"`
	FilterSpam string `json:"filter_spam,omitempty"`
}

func (c *Client) ListFwdhosts(ctx context.Context) ([]Fwdhost, error) {
	var list []Fwdhost
	err := c.getList(ctx, "/get/fwdhost/all", &list)
	return list, err
}

func (c *Client) GetFwdhost(ctx context.Context, host string) (Fwdhost, error) {
	list, err := c.ListFwdhosts(ctx)
	if err != nil {
		return Fwdhost{}, err
	}
	for _, h := range list {
		if strings.EqualFold(h.Host, host) || strings.EqualFold(h.Source, host) {
			return h, nil
		}
	}
	return Fwdhost{}, fmt.Errorf("fwdhost %q not found", host)
}

func (c *Client) CreateFwdhost(ctx context.Context, req FwdhostCreate) error {
	return c.postAction(ctx, "/add/fwdhost", req)
}

func (c *Client) DeleteFwdhost(ctx context.Context, host string) error {
	return c.postAction(ctx, "/delete/fwdhost", []string{host})
}
