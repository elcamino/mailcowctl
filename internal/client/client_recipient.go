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
