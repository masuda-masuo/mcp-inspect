// Package config parses MCP client configuration files.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ServerConfig represents a single MCP server entry.
type ServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

// Config is the parsed top-level MCP configuration.
type Config struct {
	Path    string
	Servers map[string]ServerConfig
}

// defaultPaths returns candidate config file paths in priority order.
func defaultPaths() []string {
	home, _ := os.UserHomeDir()
	return []string{
		filepath.Join(home, ".claude", "claude_desktop_config.json"),
		filepath.Join(home, ".claude.json"),
		"mcp.json",
	}
}

// Load reads a config file. If path is empty it tries the default locations.
func Load(path string) (*Config, error) {
	candidates := []string{path}
	if path == "" {
		candidates = defaultPaths()
	}

	for _, p := range candidates {
		if p == "" {
			continue
		}
		data, err := os.ReadFile(p)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("reading %s: %w", p, err)
		}
		return parse(p, data)
	}
	return nil, fmt.Errorf("no config file found; tried: %v", candidates)
}

// raw is used only for unmarshalling – supports both key spellings.
type raw struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

func parse(path string, data []byte) (*Config, error) {
	var r raw
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	if r.MCPServers == nil {
		r.MCPServers = map[string]ServerConfig{}
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	return &Config{Path: abs, Servers: r.MCPServers}, nil
}
