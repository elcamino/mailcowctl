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
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Options struct {
	Host       string
	APIKey     string
	Insecure   bool
	HTTPClient *http.Client
}

type Client struct {
	host       string
	apiKey     string
	httpClient *http.Client
}

type APIError struct {
	Type     string
	Messages []string
}

func (e *APIError) Error() string {
	if len(e.Messages) == 0 {
		return "mailcow API returned " + e.Type
	}
	return strings.Join(e.Messages, "; ")
}

type AuthError struct {
	StatusCode int
	Message    string
}

func (e *AuthError) Error() string {
	base := e.Message
	if base == "" {
		base = fmt.Sprintf("mailcow authentication failed with HTTP %d", e.StatusCode)
	}
	return base + "; check that the API key is valid, read-write for write commands, and that this source IP is allow-listed in mailcow"
}

func AsAPIError(err error, target **APIError) bool {
	return errors.As(err, target)
}

func New(opts Options) (*Client, error) {
	host := strings.TrimRight(strings.TrimSpace(opts.Host), "/")
	if host == "" {
		return nil, errors.New("host is required")
	}
	if _, err := url.ParseRequestURI(host); err != nil {
		return nil, fmt.Errorf("invalid host %q: %w", opts.Host, err)
	}
	hc := opts.HTTPClient
	if hc == nil {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		if opts.Insecure {
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
		}
		hc = &http.Client{Transport: transport}
	}
	return &Client{host: host, apiKey: opts.APIKey, httpClient: hc}, nil
}

func BoolString(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

type Domain struct {
	DomainName               string      `json:"domain_name,omitempty"`
	Domain                   string      `json:"domain,omitempty"`
	Description              string      `json:"description,omitempty"`
	Active                   any         `json:"active,omitempty"`
	MailboxesInDomain        int         `json:"mboxes_in_domain,omitempty"`
	MaxNumMailboxesForDomain int         `json:"max_num_mboxes_for_domain,omitempty"`
	QuotaUsedInDomain        string      `json:"quota_used_in_domain,omitempty"`
	MaxQuotaForDomain        json.Number `json:"max_quota_for_domain,omitempty"`
}

func (d Domain) Name() string {
	if d.DomainName != "" {
		return d.DomainName
	}
	return d.Domain
}

type Mailbox struct {
	Username  string      `json:"username,omitempty"`
	Name      string      `json:"name,omitempty"`
	Domain    string      `json:"domain,omitempty"`
	LocalPart string      `json:"local_part,omitempty"`
	Quota     json.Number `json:"quota,omitempty"`
	QuotaUsed any         `json:"quota_used,omitempty"`
	Active    any         `json:"active,omitempty"`
}

type Alias struct {
	ID      int    `json:"id"`
	Domain  string `json:"domain,omitempty"`
	Address string `json:"address,omitempty"`
	Goto    string `json:"goto,omitempty"`
	Active  any    `json:"active,omitempty"`
}

type DomainCreate struct {
	Domain             string `json:"domain"`
	Description        string `json:"description,omitempty"`
	Aliases            int    `json:"aliases,omitempty"`
	Mailboxes          int    `json:"mailboxes,omitempty"`
	DefaultQuota       int    `json:"defquota,omitempty"`
	MaxQuota           int    `json:"maxquota,omitempty"`
	Quota              int    `json:"quota,omitempty"`
	Active             string `json:"active,omitempty"`
	RateLimitValue     string `json:"rl_value,omitempty"`
	RateLimitFrame     string `json:"rl_frame,omitempty"`
	BackupMX           string `json:"backupmx,omitempty"`
	RelayAllRecipients string `json:"relay_all_recipients,omitempty"`
}

type MailboxCreate struct {
	LocalPart     string `json:"local_part"`
	Domain        string `json:"domain"`
	Name          string `json:"name"`
	Quota         string `json:"quota,omitempty"`
	Password      string `json:"password"`
	Password2     string `json:"password2"`
	Active        string `json:"active,omitempty"`
	ForcePwUpdate string `json:"force_pw_update,omitempty"`
	TLSEnforceIn  string `json:"tls_enforce_in,omitempty"`
	TLSEnforceOut string `json:"tls_enforce_out,omitempty"`
}

type AliasCreate struct {
	Address     string `json:"address"`
	Goto        string `json:"goto"`
	Active      string `json:"active,omitempty"`
	SogoVisible string `json:"sogo_visible,omitempty"`
}

func (c *Client) ListDomains(ctx context.Context) ([]Domain, error) {
	var domains []Domain
	err := c.get(ctx, "/get/domain/all", &domains)
	return domains, err
}

func (c *Client) GetDomain(ctx context.Context, domain string) (Domain, error) {
	var d Domain
	err := c.getOne(ctx, "/get/domain/"+url.PathEscape(domain), &d)
	return d, err
}

func (c *Client) CreateDomain(ctx context.Context, req DomainCreate) error {
	return c.postAction(ctx, "/add/domain", req)
}

func (c *Client) EditDomain(ctx context.Context, domain string, attr map[string]any) error {
	return c.postAction(ctx, "/edit/domain", editRequest{Attr: attr, Items: []string{domain}})
}

func (c *Client) DeleteDomain(ctx context.Context, domain string) error {
	return c.postAction(ctx, "/delete/domain", []string{domain})
}

func (c *Client) ListMailboxes(ctx context.Context, domain string) ([]Mailbox, error) {
	path := "/get/mailbox/all"
	if domain != "" {
		path += "/" + url.PathEscape(domain)
	}
	var mailboxes []Mailbox
	err := c.get(ctx, path, &mailboxes)
	return mailboxes, err
}

func (c *Client) GetMailbox(ctx context.Context, mailbox string) (Mailbox, error) {
	var m Mailbox
	err := c.getOne(ctx, "/get/mailbox/"+url.PathEscape(mailbox), &m)
	return m, err
}

func (c *Client) CreateMailbox(ctx context.Context, req MailboxCreate) error {
	return c.postAction(ctx, "/add/mailbox", req)
}

func (c *Client) EditMailbox(ctx context.Context, mailbox string, attr map[string]any) error {
	return c.postAction(ctx, "/edit/mailbox", editRequest{Attr: attr, Items: []string{mailbox}})
}

func (c *Client) DeleteMailbox(ctx context.Context, mailbox string) error {
	return c.postAction(ctx, "/delete/mailbox", []string{mailbox})
}

func (c *Client) ListAliases(ctx context.Context, domain string) ([]Alias, error) {
	var aliases []Alias
	if err := c.get(ctx, "/get/alias/all", &aliases); err != nil {
		return nil, err
	}
	if domain == "" {
		return aliases, nil
	}
	filtered := aliases[:0]
	for _, alias := range aliases {
		if alias.Domain == domain || strings.HasSuffix(alias.Address, "@"+domain) {
			filtered = append(filtered, alias)
		}
	}
	return filtered, nil
}

func (c *Client) GetAlias(ctx context.Context, ref string) (Alias, error) {
	if _, err := strconv.Atoi(ref); err == nil {
		var a Alias
		return a, c.getOne(ctx, "/get/alias/"+url.PathEscape(ref), &a)
	}
	id, err := c.ResolveAliasID(ctx, ref)
	if err != nil {
		return Alias{}, err
	}
	var a Alias
	return a, c.getOne(ctx, "/get/alias/"+strconv.Itoa(id), &a)
}

func (c *Client) CreateAlias(ctx context.Context, req AliasCreate) error {
	return c.postAction(ctx, "/add/alias", req)
}

func (c *Client) EditAlias(ctx context.Context, ref string, attr map[string]any) error {
	id, err := c.aliasID(ctx, ref)
	if err != nil {
		return err
	}
	return c.postAction(ctx, "/edit/alias", editRequest{Attr: attr, Items: []int{id}})
}

func (c *Client) DeleteAlias(ctx context.Context, ref string) error {
	id, err := c.aliasID(ctx, ref)
	if err != nil {
		return err
	}
	return c.postAction(ctx, "/delete/alias", []int{id})
}

func (c *Client) ResolveAliasID(ctx context.Context, address string) (int, error) {
	aliases, err := c.ListAliases(ctx, "")
	if err != nil {
		return 0, err
	}
	var matches []Alias
	for _, alias := range aliases {
		if strings.EqualFold(alias.Address, address) {
			matches = append(matches, alias)
		}
	}
	switch len(matches) {
	case 0:
		return 0, fmt.Errorf("alias %q not found", address)
	case 1:
		return matches[0].ID, nil
	default:
		return 0, fmt.Errorf("alias %q matched multiple ids", address)
	}
}

func (c *Client) aliasID(ctx context.Context, ref string) (int, error) {
	if id, err := strconv.Atoi(ref); err == nil {
		return id, nil
	}
	return c.ResolveAliasID(ctx, ref)
}

type editRequest struct {
	Attr  map[string]any `json:"attr"`
	Items any            `json:"items"`
}

func (c *Client) get(ctx context.Context, path string, out any) error {
	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	return dec.Decode(out)
}

func (c *Client) getOne(ctx context.Context, path string, out any) error {
	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(out); err == nil {
		return nil
	}
	var arr []json.RawMessage
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) == 0 {
		return errors.New("mailcow returned no matching object")
	}
	return json.Unmarshal(arr[0], out)
}

func (c *Client) getList(ctx context.Context, path string, out any) error {
	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("{}")) || bytes.Equal(trimmed, []byte("null")) {
		return json.Unmarshal([]byte("[]"), out)
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	return dec.Decode(out)
}

func (c *Client) postAction(ctx context.Context, path string, body any) error {
	data, err := c.do(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	return DecodeActionResponse(data)
}

func (c *Client) do(ctx context.Context, method, path string, body any) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.host+"/api/v1"+path, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mailcow transport error: %w", err)
	}
	defer resp.Body.Close()
	data, readErr := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, &AuthError{StatusCode: resp.StatusCode, Message: strings.TrimSpace(string(data))}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if readErr != nil {
			return nil, readErr
		}
		return nil, fmt.Errorf("mailcow HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	if readErr != nil {
		return nil, readErr
	}
	return data, nil
}

func DecodeActionResponse(data []byte) error {
	var actions []struct {
		Type string          `json:"type"`
		Msg  json.RawMessage `json:"msg"`
	}
	if err := json.Unmarshal(data, &actions); err != nil {
		var action struct {
			Type string          `json:"type"`
			Msg  json.RawMessage `json:"msg"`
		}
		if err2 := json.Unmarshal(data, &action); err2 != nil {
			return err
		}
		actions = append(actions, action)
	}
	for _, action := range actions {
		if action.Type != "success" {
			return &APIError{Type: action.Type, Messages: decodeMessages(action.Msg)}
		}
	}
	return nil
}

func decodeMessages(raw json.RawMessage) []string {
	if len(raw) == 0 {
		return nil
	}
	var stringsValue []string
	if err := json.Unmarshal(raw, &stringsValue); err == nil {
		return stringsValue
	}
	var single string
	if err := json.Unmarshal(raw, &single); err == nil {
		return []string{single}
	}
	var values []any
	if err := json.Unmarshal(raw, &values); err == nil {
		out := make([]string, 0, len(values))
		for _, v := range values {
			out = append(out, fmt.Sprint(v))
		}
		return out
	}
	return []string{string(raw)}
}
