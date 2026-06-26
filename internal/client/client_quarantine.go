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

type QuarantineItem struct {
	ID        int    `json:"id"`
	QID       string `json:"qid,omitempty"`
	Rcpt      string `json:"rcpt,omitempty"`
	Sender    string `json:"sender,omitempty"`
	Subject   string `json:"subject,omitempty"`
	Score     any    `json:"score,omitempty"`
	Action    string `json:"action,omitempty"`
	Created   any    `json:"created,omitempty"`
	Notified  any    `json:"notified,omitempty"`
	VirusFlag any    `json:"virus_flag,omitempty"`
}

func (c *Client) ListQuarantine(ctx context.Context, rcpt string) ([]QuarantineItem, error) {
	list, err := apiList[QuarantineItem](ctx, c, "/get/quarantine/all")
	if err != nil {
		return nil, err
	}
	if rcpt == "" {
		return list, nil
	}
	filtered := list[:0]
	for _, q := range list {
		if strings.EqualFold(q.Rcpt, rcpt) {
			filtered = append(filtered, q)
		}
	}
	return filtered, nil
}

func (c *Client) GetQuarantine(ctx context.Context, id int) (QuarantineItem, error) {
	list, err := c.ListQuarantine(ctx, "")
	if err != nil {
		return QuarantineItem{}, err
	}
	return findByID(list, id, func(q QuarantineItem) int { return q.ID }, "quarantine item")
}

func (c *Client) quarantineAction(ctx context.Context, id int, action string) error {
	return c.postAction(ctx, "/edit/qitem", editRequest{Attr: map[string]any{"action": action}, Items: []int{id}})
}

func (c *Client) ReleaseQuarantine(ctx context.Context, id int) error {
	return c.quarantineAction(ctx, id, "release")
}

func (c *Client) LearnHamQuarantine(ctx context.Context, id int) error {
	return c.quarantineAction(ctx, id, "learnham")
}

func (c *Client) DeleteQuarantine(ctx context.Context, id int) error {
	body := struct {
		Items []int `json:"items"`
	}{Items: []int{id}}
	return c.postAction(ctx, "/delete/qitem", body)
}
