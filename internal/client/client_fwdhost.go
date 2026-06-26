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
	return apiList[Fwdhost](ctx, c, "/get/fwdhost/all")
}

func (c *Client) GetFwdhost(ctx context.Context, host string) (Fwdhost, error) {
	list, err := c.ListFwdhosts(ctx)
	if err != nil {
		return Fwdhost{}, err
	}
	return findFirst(list, func(h Fwdhost) bool {
		return strings.EqualFold(h.Host, host) || strings.EqualFold(h.Source, host)
	}, fmt.Errorf("fwdhost %q not found", host))
}

func (c *Client) CreateFwdhost(ctx context.Context, req FwdhostCreate) error {
	return c.postAction(ctx, "/add/fwdhost", req)
}

func (c *Client) DeleteFwdhost(ctx context.Context, host string) error {
	return c.postAction(ctx, "/delete/fwdhost", []string{host})
}
