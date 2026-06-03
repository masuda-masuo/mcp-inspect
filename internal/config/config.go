// Package config parses MCP client configuration files.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
// On Windows, Claude Desktop is installed as a UWP app whose data directory
// is under AppData\Local\Packages\Claude_<random>\LocalCache\Roaming\Claude\.
// We glob for that pattern so users don't need to specify --config manually.
func defaultPaths() []string {
	home, _ := os.UserHomeDir()

	candidates := []string{
		// Standard / cross-platform locations
		filepath.Join(home, ".claude", "claude_desktop_config.json"),
		filepath.Join(home, ".claude.json"),
		"mcp.json",
	}

	if runtime.GOOS == "windows" {
		candidates = append(windowsClaudePaths(home), candidates...)
	}

	return candidates
}

// windowsClaudePaths globs for Claude Desktop's UWP data directory.
// The package folder name has the form "Claude_<id>" where <id> varies
// per installation, so we cannot hard-code the full path.
func windowsClaudePaths(home string) []string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = filepath.Join(home, "AppData", "Local")
	}

	pattern := filepath.Join(
		localAppData,
		"Packages", "Claude_*",
		"LocalCache", "Roaming", "Claude",
		"claude_desktop_config.json",
	)

	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return nil
	}
	return matches
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
