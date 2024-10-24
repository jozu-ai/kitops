// Copyright 2024 The KitOps Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0
package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"kitops/pkg/output"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Config struct {
	LogLevel string `json:"logLevel"`
	Progress string `json:"progress"`
}

// DefaultConfig returns a Config struct with default values.
func DefaultConfig() *Config {
	return &Config{
		LogLevel: output.LogLevelInfo.String(),
		Progress: "plain",
	}
}

// Set a configuration key and value.
func setConfig(_ context.Context, opts *configOptions) error {
	configPath := getConfigPath(opts.profile)
	cfg, err := LoadConfig(configPath)
	if err != nil {
		cfg = DefaultConfig()
	}

	v := reflect.ValueOf(cfg).Elem().FieldByName(cases.Title(language.Und).String(opts.key))
	if !v.IsValid() {
		return fmt.Errorf("unknown configuration key: %s", opts.key)
	}

	v.SetString(opts.value)
	err = SaveConfig(cfg, configPath) // Save only when configuration is modified.
	if err != nil {
		return err
	}
	fmt.Printf("Config '%s' set to '%s'\n", opts.key, opts.value)
	return nil
}

// Get a configuration value.
func getConfig(_ context.Context, opts *configOptions) (string, error) {
	configPath := getConfigPath(opts.profile)
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return "", err
	}

	v := reflect.ValueOf(cfg).Elem().FieldByName(strings.Title(opts.key))
	if !v.IsValid() {
		return "", fmt.Errorf("unknown configuration key: %s", opts.key)
	}

	return fmt.Sprintf("%v", v.Interface()), nil
}

// List all configuration values.
func listConfig(_ context.Context, opts *configOptions) error {
	configPath := getConfigPath(opts.profile)
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	// Use reflection to iterate through fields and print them.
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("%s: %v\n", t.Field(i).Name, v.Field(i).Interface())
	}
	return nil
}

// Reset configuration to defaults.
func resetConfig(_ context.Context, opts *configOptions) error {
	configPath := getConfigPath(opts.profile)
	cfg := DefaultConfig()
	err := SaveConfig(cfg, configPath)
	if err != nil {
		return err
	}
	fmt.Println("Configuration reset to default values.")
	return nil
}

// Load configuration from a file.
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	file, err := os.Open(configPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// Return default config, but don't save it to the file.
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}

	// If some fields are empty, fallback to defaults.
	defaultConfig := DefaultConfig()
	if config.LogLevel == "" {
		config.LogLevel = defaultConfig.LogLevel
	}
	if config.Progress == "" {
		config.Progress = defaultConfig.Progress
	}

	return &config, nil
}

// Save configuration to a file.
func SaveConfig(config *Config, configPath string) error {
	// Ensure the directory exists before saving the file.
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(config)
}

// Get the config path, either from the profile or default.
func getConfigPath(profile string) string {
	configDir := os.Getenv("KITOPS_HOME")
	if configDir == "" {
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".kitops")
	}
	if profile != "" {
		configDir = filepath.Join(configDir, "profiles", profile)
	}
	return filepath.Join(configDir, "config.json")
}

// ConfigOptions struct to store command options.
type configOptions struct {
	key        string
	value      string
	profile    string
	configHome string
}
