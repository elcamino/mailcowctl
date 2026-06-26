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
	list, err := apiList[Filter](ctx, c, "/get/filters/all")
	if err != nil {
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
	return findByID(list, id, func(f Filter) int { return f.ID }, "filter")
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
