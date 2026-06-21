package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Inputs struct {
	Host       string
	APIKey     string
	Profile    string
	ConfigPath string
}

type Config struct {
	Host    string
	APIKey  string
	Profile string
}

type fileConfig struct {
	CurrentProfile string                   `yaml:"current_profile"`
	Profiles       map[string]profileConfig `yaml:"profiles"`
}

type profileConfig struct {
	Host   string `yaml:"host"`
	APIKey string `yaml:"api_key"`
}

type AuthConfigError struct {
	Field string
}

func (e *AuthConfigError) Error() string {
	return fmt.Sprintf("missing mailcow %s; set --%s, MAILCOW_%s, or a config profile", e.Field, flagName(e.Field), envName(e.Field))
}

func IsAuthConfigError(err error) bool {
	var target *AuthConfigError
	return errors.As(err, &target)
}

func Resolve(in Inputs) (Config, error) {
	fc, _ := readFileConfig(configPath(in.ConfigPath))

	profile := firstNonEmpty(in.Profile, os.Getenv("MAILCOW_PROFILE"), fc.CurrentProfile)
	cfg := Config{Profile: profile}
	if profile != "" && fc.Profiles != nil {
		if p, ok := fc.Profiles[profile]; ok {
			cfg.Host = p.Host
			cfg.APIKey = p.APIKey
		}
	}

	if v := os.Getenv("MAILCOW_HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("MAILCOW_API_KEY"); v != "" {
		cfg.APIKey = v
	}
	if in.Host != "" {
		cfg.Host = in.Host
	}
	if in.APIKey != "" {
		cfg.APIKey = in.APIKey
	}
	if in.Profile != "" {
		cfg.Profile = in.Profile
	}
	cfg.Host = NormalizeHost(cfg.Host)

	if cfg.Host == "" {
		return Config{}, &AuthConfigError{Field: "host"}
	}
	if cfg.APIKey == "" {
		return Config{}, &AuthConfigError{Field: "api key"}
	}
	return cfg, nil
}

func CurrentProfile(path string) (string, error) {
	fc, err := readFileConfig(configPath(path))
	if err != nil {
		return "", err
	}
	return fc.CurrentProfile, nil
}

func SetCurrentProfile(path, name string) error {
	if name == "" {
		return errors.New("profile name is required")
	}
	p := configPath(path)
	fc, _ := readFileConfig(p)
	if fc.Profiles == nil {
		fc.Profiles = map[string]profileConfig{}
	}
	fc.CurrentProfile = name
	data, err := yaml.Marshal(&fc)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o600)
}

func NormalizeHost(host string) string {
	return strings.TrimRight(strings.TrimSpace(host), "/")
}

func ConfigPath(path string) string {
	return configPath(path)
}

func configPath(path string) string {
	if path != "" {
		return path
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "mailcowctl", "config.yaml")
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".config", "mailcowctl", "config.yaml")
	}
	return filepath.Join(".config", "mailcowctl", "config.yaml")
}

func readFileConfig(path string) (fileConfig, error) {
	var fc fileConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return fc, err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return fc, nil
	}
	if err := yaml.Unmarshal(data, &fc); err != nil {
		return fc, err
	}
	return fc, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func flagName(field string) string {
	return strings.ReplaceAll(field, " ", "-")
}

func envName(field string) string {
	return strings.ToUpper(strings.ReplaceAll(field, " ", "_"))
}
