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

func apiList[T any](ctx context.Context, c *Client, path string) ([]T, error) {
	var list []T
	err := c.getList(ctx, path, &list)
	return list, err
}

func findFirst[T any](items []T, match func(T) bool, notFound error) (T, error) {
	for _, item := range items {
		if match(item) {
			return item, nil
		}
	}
	var zero T
	return zero, notFound
}

func findByID[T any](items []T, id int, idOf func(T) int, label string) (T, error) {
	return findFirst(items, func(item T) bool {
		return idOf(item) == id
	}, fmt.Errorf("%s %d not found", label, id))
}
