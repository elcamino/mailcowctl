// SPDX-License-Identifier: GPL-3.0-or-later
//
// Copyright (C) 2026 Tobias von Dewitz
//
// This file is part of mailcowctl.
//
// mailcowctl is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// mailcowctl is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with mailcowctl. If not, see <https://www.gnu.org/licenses/>.

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
