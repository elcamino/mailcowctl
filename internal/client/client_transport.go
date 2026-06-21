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

type Transport struct {
	ID            int    `json:"id"`
	Destination   string `json:"destination,omitempty"`
	Nexthop       string `json:"nexthop,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	PasswordShort string `json:"password_short,omitempty"`
	Active        any    `json:"active,omitempty"`
}

type TransportCreate struct {
	Destination string `json:"destination"`
	Nexthop     string `json:"nexthop"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Active      string `json:"active,omitempty"`
}

func (c *Client) ListTransports(ctx context.Context) ([]Transport, error) {
	var list []Transport
	err := c.getList(ctx, "/get/transport/all", &list)
	return list, err
}

func (c *Client) GetTransport(ctx context.Context, id int) (Transport, error) {
	list, err := c.ListTransports(ctx)
	if err != nil {
		return Transport{}, err
	}
	for _, t := range list {
		if t.ID == id {
			return t, nil
		}
	}
	return Transport{}, fmt.Errorf("transport %d not found", id)
}

func (c *Client) CreateTransport(ctx context.Context, req TransportCreate) error {
	return c.postAction(ctx, "/add/transport", req)
}

func (c *Client) DeleteTransport(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/transport", []int{id})
}
