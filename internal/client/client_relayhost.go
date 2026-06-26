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
)

type Relayhost struct {
	ID              int    `json:"id"`
	Hostname        string `json:"hostname,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	PasswordShort   string `json:"password_short,omitempty"`
	Active          any    `json:"active,omitempty"`
	UsedByDomains   string `json:"used_by_domains,omitempty"`
	UsedByMailboxes string `json:"used_by_mailboxes,omitempty"`
}

type RelayhostCreate struct {
	Hostname string `json:"hostname"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Active   string `json:"active,omitempty"`
}

func (c *Client) ListRelayhosts(ctx context.Context) ([]Relayhost, error) {
	return apiList[Relayhost](ctx, c, "/get/relayhost/all")
}

func (c *Client) GetRelayhost(ctx context.Context, id int) (Relayhost, error) {
	list, err := c.ListRelayhosts(ctx)
	if err != nil {
		return Relayhost{}, err
	}
	return findByID(list, id, func(h Relayhost) int { return h.ID }, "relayhost")
}

func (c *Client) CreateRelayhost(ctx context.Context, req RelayhostCreate) error {
	return c.postAction(ctx, "/add/relayhost", req)
}

func (c *Client) DeleteRelayhost(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/relayhost", []int{id})
}
