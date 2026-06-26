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

type Resource struct {
	ID               int    `json:"id"`
	Description      string `json:"description,omitempty"`
	Domain           string `json:"domain,omitempty"`
	Kind             string `json:"kind,omitempty"`
	MultipleBookings any    `json:"multiple_bookings,omitempty"`
	Active           any    `json:"active,omitempty"`
}

type ResourceCreate struct {
	Description      string `json:"description"`
	Domain           string `json:"domain"`
	Kind             string `json:"kind"`
	MultipleBookings int    `json:"multiple_bookings"`
	Active           string `json:"active,omitempty"`
}

func (c *Client) ListResources(ctx context.Context) ([]Resource, error) {
	return apiList[Resource](ctx, c, "/get/resource/all")
}

func (c *Client) GetResource(ctx context.Context, id int) (Resource, error) {
	list, err := c.ListResources(ctx)
	if err != nil {
		return Resource{}, err
	}
	return findByID(list, id, func(r Resource) int { return r.ID }, "resource")
}

func (c *Client) CreateResource(ctx context.Context, req ResourceCreate) error {
	return c.postAction(ctx, "/add/resource", req)
}

func (c *Client) DeleteResource(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/resource", []int{id})
}
