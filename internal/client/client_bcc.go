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

type Bcc struct {
	ID        int    `json:"id"`
	Type      string `json:"type,omitempty"`
	LocalDest string `json:"local_dest,omitempty"`
	BccDest   string `json:"bcc_dest,omitempty"`
	Domain    string `json:"domain,omitempty"`
	Active    any    `json:"active,omitempty"`
}

type BccCreate struct {
	LocalDest string `json:"local_dest"`
	BccDest   string `json:"bcc_dest"`
	Type      string `json:"type"`
	Active    string `json:"active,omitempty"`
}

func (c *Client) ListBccs(ctx context.Context) ([]Bcc, error) {
	var list []Bcc
	err := c.getList(ctx, "/get/bcc/all", &list)
	return list, err
}

func (c *Client) GetBcc(ctx context.Context, id int) (Bcc, error) {
	list, err := c.ListBccs(ctx)
	if err != nil {
		return Bcc{}, err
	}
	for _, b := range list {
		if b.ID == id {
			return b, nil
		}
	}
	return Bcc{}, fmt.Errorf("bcc %d not found", id)
}

func (c *Client) CreateBcc(ctx context.Context, req BccCreate) error {
	return c.postAction(ctx, "/add/bcc", req)
}

func (c *Client) DeleteBcc(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/bcc", []int{id})
}
