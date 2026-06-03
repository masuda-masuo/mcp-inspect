// mcp-inspect generates a human-readable report of MCP tools from a config file.
//
// Usage:
//
//	mcp-inspect [flags]
//
// Flags:
//
//	--config PATH   Path to MCP config file (default: auto-discover)
//	--output PATH   Output file path (default: ./mcp-report.html)
//	--format        Output format: html (default) or json
//	--no-open       Do not open the report in a browser after generation
//	--no-launch     Do not launch servers; show config contents only
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/masuda-masuo/mcp-inspect/internal/config"
	"github.com/masuda-masuo/mcp-inspect/internal/fetcher"
	"github.com/masuda-masuo/mcp-inspect/internal/i18n"
	"github.com/masuda-masuo/mcp-inspect/internal/reporter"
)

var version = "dev"

func main() {
	configPath := flag.String("config", "", "Path to MCP config file")
	outputPath := flag.String("output", "", "Output file path (default: ./mcp-report.html or stdout for json)")
	format := flag.String("format", "html", "Output format: html or json")
	noOpen := flag.Bool("no-open", false, "Do not open the report in the browser")
	noLaunch := flag.Bool("no-launch", false, "Do not launch servers; show config contents only (no tool list)")
	lang := flag.String("lang", "en", "Output language: en (default) or ja")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println("mcp-inspect", version)
		return
	}

	if err := run(*configPath, *outputPath, *format, *noOpen, *noLaunch, i18n.Parse(*lang)); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(configPath, outputPath, format string, noOpen, noLaunch bool, lang i18n.Lang) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	fmt.Fprintf(os.Stderr, "\u2192 Config: %s (%d servers)\n", cfg.Path, len(cfg.Servers))

	if len(cfg.Servers) == 0 {
		fmt.Fprintln(os.Stderr, "  No servers found in config.")
	}

	var report *reporter.Report
	if noLaunch {
		fmt.Fprintln(os.Stderr, "\u2192 --no-launch: skipping server startup, showing config only")
		report = reporter.BuildNoLaunch(cfg, lang)
	} else {
		fmt.Fprintln(os.Stderr, "\u2192 Launching servers and fetching tool lists...")
		results := fetcher.FetchAll(cfg)
		for _, r := range results {
			if r.Error != "" {
				fmt.Fprintf(os.Stderr, "  [%s] \u26a0 error: %s\n", r.ServerName, r.Error)
			} else {
				fmt.Fprintf(os.Stderr, "  [%s] %d tools\n", r.ServerName, len(r.Tools))
			}
		}
		report = reporter.Build(cfg, results, lang)
	}

	switch strings.ToLower(format) {
	case "json":
		w := os.Stdout
		if outputPath != "" {
			f, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			defer f.Close()
			w = f
			fmt.Fprintf(os.Stderr, "\u2192 JSON written to %s\n", outputPath)
		}
		return reporter.WriteJSON(w, report)

	case "html", "":
		if outputPath == "" {
			outputPath = "mcp-report.html"
		}
		f, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := reporter.WriteHTML(f, report); err != nil {
			return err
		}
		abs, _ := filepath.Abs(outputPath)
		fmt.Fprintf(os.Stderr, "\u2192 Report written to %s\n", abs)
		if !noOpen {
			openBrowser(abs)
		}
		return nil

	default:
		return fmt.Errorf("unknown format %q (use html or json)", format)
	}
}

func openBrowser(path string) {
	url := "file://" + path
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
