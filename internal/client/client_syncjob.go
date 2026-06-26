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

type SyncJob struct {
	ID           int    `json:"id"`
	Mailbox      string `json:"user2,omitempty"`
	Host1        string `json:"host1,omitempty"`
	Port1        any    `json:"port1,omitempty"`
	User1        string `json:"user1,omitempty"`
	Enc1         string `json:"enc1,omitempty"`
	MinsInterval any    `json:"mins_interval,omitempty"`
	Active       any    `json:"active,omitempty"`
	LastRun      string `json:"last_run,omitempty"`
	IsRunning    any    `json:"is_running,omitempty"`
}

type SyncJobCreate struct {
	Username          string `json:"username"`
	Host1             string `json:"host1"`
	Port1             int    `json:"port1"`
	User1             string `json:"user1"`
	Password1         string `json:"password1"`
	Enc1              string `json:"enc1,omitempty"`
	MinsInterval      int    `json:"mins_interval,omitempty"`
	MaxAge            int    `json:"maxage,omitempty"`
	Active            string `json:"active,omitempty"`
	Delete1           string `json:"delete1,omitempty"`
	Delete2           string `json:"delete2,omitempty"`
	Delete2Duplicates string `json:"delete2duplicates,omitempty"`
	SubscribeAll      string `json:"subscribeall,omitempty"`
	Exclude           string `json:"exclude,omitempty"`
	Automap           string `json:"automap,omitempty"`
}

func (c *Client) ListSyncJobs(ctx context.Context, mailbox string) ([]SyncJob, error) {
	jobs, err := apiList[SyncJob](ctx, c, "/get/syncjobs/all/no_log")
	if err != nil {
		return nil, err
	}
	if mailbox == "" {
		return jobs, nil
	}
	filtered := jobs[:0]
	for _, j := range jobs {
		if strings.EqualFold(j.Mailbox, mailbox) {
			filtered = append(filtered, j)
		}
	}
	return filtered, nil
}

func (c *Client) GetSyncJob(ctx context.Context, id int) (SyncJob, error) {
	jobs, err := c.ListSyncJobs(ctx, "")
	if err != nil {
		return SyncJob{}, err
	}
	return findByID(jobs, id, func(j SyncJob) int { return j.ID }, "syncjob")
}

func (c *Client) CreateSyncJob(ctx context.Context, req SyncJobCreate) error {
	return c.postAction(ctx, "/add/syncjob", req)
}

func (c *Client) EditSyncJob(ctx context.Context, id int, attr map[string]any) error {
	return c.postAction(ctx, "/edit/syncjob", editRequest{Attr: attr, Items: []int{id}})
}

func (c *Client) DeleteSyncJob(ctx context.Context, id int) error {
	return c.postAction(ctx, "/delete/syncjob", []int{id})
}
