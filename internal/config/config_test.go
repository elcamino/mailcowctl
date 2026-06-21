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

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveUsesFlagsThenEnvThenProfileFile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	configDir := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "mailcowctl")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(`
current_profile: prod
profiles:
  prod:
    host: https://file.example
    api_key: file-key
  staging:
    host: https://staging.example
    api_key: staging-key
`), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("MAILCOW_HOST", "https://env.example")
	t.Setenv("MAILCOW_API_KEY", "env-key")
	t.Setenv("MAILCOW_PROFILE", "staging")

	cfg, err := Resolve(Inputs{Host: "https://flag.example/"})
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if cfg.Host != "https://flag.example" {
		t.Fatalf("host = %q, want flag host normalized without trailing slash", cfg.Host)
	}
	if cfg.APIKey != "env-key" {
		t.Fatalf("api key = %q, want env key", cfg.APIKey)
	}
	if cfg.Profile != "staging" {
		t.Fatalf("profile = %q, want env profile", cfg.Profile)
	}
}

func TestResolveFallsBackToCurrentProfile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	configDir := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "mailcowctl")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(`
current_profile: prod
profiles:
  prod:
    host: https://file.example/
    api_key: file-key
`), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Resolve(Inputs{})
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if cfg.Host != "https://file.example" || cfg.APIKey != "file-key" || cfg.Profile != "prod" {
		t.Fatalf("cfg = %+v, want current profile values", cfg)
	}
}

func TestResolveMissingAuthReturnsConfigError(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	_, err := Resolve(Inputs{Host: "https://mail.example"})
	if err == nil {
		t.Fatal("Resolve returned nil error, want missing API key error")
	}
	if !IsAuthConfigError(err) {
		t.Fatalf("IsAuthConfigError(%v) = false, want true", err)
	}
}
