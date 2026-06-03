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
func defaultPaths() []string {
	home, _ := os.UserHomeDir()

	var candidates []string

	switch runtime.GOOS {
	case "windows":
		// Claude Desktop on Windows ships as a UWP app.
		// Primary path: %APPDATA%\Claude\claude_desktop_config.json
		// Fallback: glob the UWP package directory (package id suffix varies per install).
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		candidates = append(candidates,
			filepath.Join(appData, "Claude", "claude_desktop_config.json"),
		)
		candidates = append(candidates, windowsUWPPaths(home)...)

	case "darwin":
		candidates = append(candidates,
			filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json"),
		)

	default: // linux and others
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			configHome = filepath.Join(home, ".config")
		}
		candidates = append(candidates,
			filepath.Join(configHome, "Claude", "claude_desktop_config.json"),
		)
	}

	// Cross-platform fallbacks (Claude Code / manual locations)
	candidates = append(candidates,
		filepath.Join(home, ".claude", "claude_desktop_config.json"),
		filepath.Join(home, ".claude.json"),
		"mcp.json",
	)

	return candidates
}

// windowsUWPPaths globs for the Claude Desktop UWP package directory.
// The folder name has the form "Claude_<random-id>" and varies per installation.
func windowsUWPPaths(home string) []string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = filepath.Join(home, "AppData", "Local")
	}
	pattern := filepath.Join(
		localAppData, "Packages", "Claude_*",
		"LocalCache", "Roaming", "Claude", "claude_desktop_config.json",
	)
	matches, err := filepath.Glob(pattern)
	if err != nil {
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

// raw is used only for unmarshalling.
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
