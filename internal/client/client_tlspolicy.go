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
	return apiList[TlsPolicy](ctx, c, "/get/tls-policy-map/all")
}

func (c *Client) GetTlsPolicy(ctx context.Context, id int) (TlsPolicy, error) {
	list, err := c.ListTlsPolicies(ctx)
	if err != nil {
		return TlsPolicy{}, err
	}
	return findByID(list, id, func(p TlsPolicy) int { return p.ID }, "tls policy")
}

func (c *Client) CreateTlsPolicy(ctx context.Context, req TlsPolicyCreate) error {
	return c.postAction(ctx, "/add/tls-policy-map", req)
}

func (c *Client) DeleteTlsPolicy(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/tls-policy-map", []int{id})
}
