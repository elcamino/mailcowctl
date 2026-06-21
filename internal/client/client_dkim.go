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

type Dkim struct {
	Domain   string `json:"-"`
	Selector string `json:"dkim_selector,omitempty"`
	Length   string `json:"length,omitempty"`
	PrivKey  string `json:"privkey,omitempty"` // always empty: the mailcow API never returns DKIM private keys
	PubKey   string `json:"pubkey,omitempty"`
	DkimTxt  string `json:"dkim_txt,omitempty"`
}

type DkimCreate struct {
	Domains  string `json:"domains"`
	Selector string `json:"dkim_selector"`
	KeySize  int    `json:"key_size"`
}

type DkimDuplicate struct {
	FromDomain string `json:"from_domain"`
	ToDomain   string `json:"to_domain"`
}

func (c *Client) GetDkim(ctx context.Context, domain string) (Dkim, error) {
	var d Dkim
	if err := c.getOne(ctx, "/get/dkim/"+url.PathEscape(domain), &d); err != nil {
		return Dkim{}, err
	}
	d.Domain = domain
	return d, nil
}

func (c *Client) CreateDkim(ctx context.Context, req DkimCreate) error {
	return c.postAction(ctx, "/add/dkim", req)
}

func (c *Client) DuplicateDkim(ctx context.Context, req DkimDuplicate) error {
	return c.postAction(ctx, "/add/dkim_duplicate", req)
}

func (c *Client) DeleteDkim(ctx context.Context, domain string) error {
	return c.postAction(ctx, "/delete/dkim", []string{domain})
}
