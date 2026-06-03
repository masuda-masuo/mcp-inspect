// Package reporter provides HTML and JSON report generation.
package reporter

import (
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/masuda-masuo/mcp-inspect/internal/classifier"
	"github.com/masuda-masuo/mcp-inspect/internal/i18n"
	"github.com/masuda-masuo/mcp-inspect/internal/config"
	"github.com/masuda-masuo/mcp-inspect/internal/fetcher"
)

// ToolReport is a tool enriched with warning badges.
type ToolReport struct {
	Name        string
	Description string
	Warnings    []classifier.Warning
}

// ServerReport aggregates one server's inspection result.
type ServerReport struct {
	Name      string
	Command   string
	Args      []string
	EnvKeys   []string // env var names only (not values, for privacy)
	Tools     []ToolReport
	Error     string
	ToolCount int
	WarnCount int
	NoLaunch  bool // true when built from config only (--no-launch)

	// Enrichment from config (available in both normal and no-launch modes)
	RuntimeKind  string   // "npx" | "node" | "uvx" | "python" | "binary"
	PackageName  string   // e.g. "@modelcontextprotocol/server-filesystem"
	AllowedPaths []string // filesystem-style servers: list of path args
}

// Report is the top-level report structure.
type Report struct {
	GeneratedAt   time.Time
	ConfigPath    string
	Servers       []ServerReport
	TotalTools    int
	TotalWarnings int
	WarnBreakdown map[classifier.Warning]int
	NoLaunch      bool
	Lang          i18n.Lang
}

// Build creates a Report from fetcher results + config (normal mode).
func Build(cfg *config.Config, results []fetcher.Result, lang i18n.Lang) *Report {
	r := &Report{
		GeneratedAt:   time.Now().UTC(),
		ConfigPath:    cfg.Path,
		WarnBreakdown: map[classifier.Warning]int{},
		Lang:          lang,
	}

	for _, res := range results {
		sr := ServerReport{
			Name:    res.ServerName,
			Command: res.Command,
			Error:   res.Error,
		}

		// Enrich with config info (args, env, runtime kind)
		if srv, ok := cfg.Servers[res.ServerName]; ok {
			envKeys := make([]string, 0, len(srv.Env))
			for k := range srv.Env {
				envKeys = append(envKeys, k)
			}
			sort.Strings(envKeys)
			sr.Args = srv.Args
			sr.EnvKeys = envKeys
			enrichNoLaunch(&sr)
		}

		for _, t := range res.Tools {
			warns := classifier.Classify(t.Name, t.Description)
			tr := ToolReport{
				Name:        t.Name,
				Description: t.Description,
				Warnings:    warns,
			}
			sr.Tools = append(sr.Tools, tr)
			if len(warns) > 0 {
				sr.WarnCount++
			}
			for _, w := range warns {
				r.WarnBreakdown[w]++
			}
		}
		sr.ToolCount = len(sr.Tools)
		r.TotalTools += sr.ToolCount
		r.TotalWarnings += sr.WarnCount
		r.Servers = append(r.Servers, sr)
	}
	return r
}

// BuildNoLaunch creates a Report from config only, without launching servers.
func BuildNoLaunch(cfg *config.Config, lang i18n.Lang) *Report {
	r := &Report{
		GeneratedAt:   time.Now().UTC(),
		ConfigPath:    cfg.Path,
		WarnBreakdown: map[classifier.Warning]int{},
		NoLaunch:      true,
		Lang:          lang,
	}

	for name, srv := range cfg.Servers {
		envKeys := make([]string, 0, len(srv.Env))
		for k := range srv.Env {
			envKeys = append(envKeys, k)
		}
		sort.Strings(envKeys)

		sr := ServerReport{
			Name:     name,
			Command:  srv.Command,
			Args:     srv.Args,
			EnvKeys:  envKeys,
			NoLaunch: true,
		}
		enrichNoLaunch(&sr)
		r.Servers = append(r.Servers, sr)
	}
	return r
}

// enrichNoLaunch fills RuntimeKind, PackageName, AllowedPaths from command+args.
func enrichNoLaunch(sr *ServerReport) {
	cmd := strings.ToLower(filepath.Base(sr.Command))
	cmd = strings.TrimSuffix(cmd, ".exe")

	switch cmd {
	case "npx", "bunx":
		sr.RuntimeKind = "npx"
		for _, a := range sr.Args {
			if !strings.HasPrefix(a, "-") {
				sr.PackageName = a
				break
			}
		}
	case "node":
		sr.RuntimeKind = "node"
		if len(sr.Args) > 0 {
			sr.PackageName = nodeScriptToPackage(sr.Args[0])
		}
		if len(sr.Args) > 1 {
			sr.AllowedPaths = pathArgs(sr.Args[1:])
		}
	case "uvx", "uv":
		sr.RuntimeKind = "uvx"
		for _, a := range sr.Args {
			if !strings.HasPrefix(a, "-") {
				sr.PackageName = a
				break
			}
		}
	case "python", "python3":
		sr.RuntimeKind = "python"
		for i, a := range sr.Args {
			if a == "-m" && i+1 < len(sr.Args) {
				sr.PackageName = sr.Args[i+1]
				break
			}
			if !strings.HasPrefix(a, "-") {
				sr.PackageName = filepath.Base(a)
				break
			}
		}
	default:
		sr.RuntimeKind = "binary"
		// binary: package name not needed, command itself is shown
		sr.AllowedPaths = pathArgs(sr.Args)
	}
}

func nodeScriptToPackage(scriptPath string) string {
	p := strings.ReplaceAll(scriptPath, "\\", "/")
	idx := strings.Index(p, "node_modules/")
	if idx < 0 {
		return filepath.Base(scriptPath)
	}
	rest := p[idx+len("node_modules/"):]
	parts := strings.SplitN(rest, "/", 3)
	if len(parts) >= 2 && strings.HasPrefix(parts[0], "@") {
		return parts[0] + "/" + parts[1]
	}
	return parts[0]
}

func pathArgs(args []string) []string {
	var paths []string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") && looksLikePath(a) {
			paths = append(paths, a)
		}
	}
	return paths
}

func looksLikePath(s string) bool {
	return strings.ContainsAny(s, "/\\") || strings.HasPrefix(s, "~")
}
