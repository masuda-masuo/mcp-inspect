# mcp-inspect

> Visualize your MCP tool surface before you hand the keys to an AI.

**"Read it once before you let the AI use it."**

`mcp-inspect` reads your MCP client config, launches each server, and generates an HTML report of every tool it exposes — with badges highlighting destructive, write, or outbound-send operations.

---

## The idea

```
Using tools you don't understand → AI can do things you didn't expect → this is the real risk
```

Two things `mcp-inspect` gives you:

1. **The actual tool list, side-by-side with the README you already have.**
   Run `tools/list` and render the result. You spot "the README says 3 tools, but 6 showed up."
   Diffing against the README is intentionally left to the human eye.

2. **Tools that look risky, surfaced automatically.**
   Names like `delete_repo`, `transfer_repo`, `send_report` get a badge.
   Keyword matching only — accuracy is traded for simplicity.
   The goal is *awareness*, not a security gate.

```
Your README (open in another tab):    mcp-inspect report:
  - create_issue                        - create_issue
  - list_issues                         - list_issues
  - get_pull_request                    - get_pull_request
                                        - delete_repo        ← destructive
                                        - add_collaborator
                                        - transfer_repo      ← destructive
```

`mcp-inspect` does not block anything. Blocking is [mcp-launcher](https://github.com/masuda-masuo/mcp-launcher)'s job.

---

## Install

Download the binary for your platform from the [latest release](https://github.com/masuda-masuo/mcp-inspect/releases/latest):

| Platform | File |
|----------|------|
| Windows (amd64) | `mcp-inspect-windows-amd64.exe` |
| macOS (Apple Silicon) | `mcp-inspect-darwin-arm64` |
| macOS (Intel) | `mcp-inspect-darwin-amd64` |
| Linux (amd64) | `mcp-inspect-linux-amd64` |

On macOS/Linux, make it executable: `chmod +x mcp-inspect-*`

---

## Usage

```
mcp-inspect [flags]

Flags:
  --config PATH    Path to MCP config file (default: auto-discover)
  --output PATH    Output file (default: ./mcp-report.html)
  --format         Output format: html (default) or json
  --no-open        Do not open the report in the browser after generation
  --no-launch      Do not start servers; show config contents only
  --version        Print version and exit
```

### Default config discovery order

1. `~/.claude/claude_desktop_config.json`
2. `~/.claude.json`
3. `mcp.json` in the current directory

### Typical workflow

```bash
# Inspect your Claude Desktop config (opens browser automatically)
mcp-inspect

# Explicit config path
mcp-inspect --config ~/.claude/claude_desktop_config.json

# JSON output for CI
mcp-inspect --format json | jq .

# Config-only view without starting any servers
mcp-inspect --no-launch
```

---

## HTML report

Each server card shows:

- **Tool table** with name, description, and warning badges
- **Config panel** (collapsible) — runtime type, package name, allowed directories, env var names

### Badges

Badges are awareness prompts, not security verdicts. Colors are deliberately muted.

| Badge | Triggered by (tool name keywords) |
|-------|-----------------------------------|
| Destructive | `delete` `remove` `destroy` `drop` `purge` `transfer` `truncate` `wipe` … |
| Write | `write` `update` `push` `create` `edit` `modify` `upload` `commit` … |
| External Send | `send` `publish` `notify` `webhook` `email` `tweet` … |

### JSON output

```bash
mcp-inspect --format json
```

```json
{
  "generated_at": "2026-06-03T12:00:00Z",
  "config": "~/.claude/claude_desktop_config.json",
  "servers": [
    {
      "name": "github",
      "command": "github-mcp-server",
      "tool_count": 12,
      "tools": [
        { "name": "delete_repo", "description": "...", "warnings": ["destructive"] }
      ]
    }
  ]
}
```

---

## `--no-launch` mode

Reads only the config file — no server processes are started.

Useful for a quick inventory of what's registered, without any execution risk.
Each server card shows command, args, allowed directories (for filesystem-style servers), and env var names (values are never shown).

> **Note:** Starting a server to read its tool list is the normal risk model.
> An MCP stdio server can run arbitrary code the moment its process starts — before `initialize` is even sent.
> This is an inherent property of the protocol, not a bug in `mcp-inspect`.
> For servers you don't trust, don't use them at all.

---

## Demo (testdata)

```bash
# Build mock servers and run the demo
make demo          # macOS/Linux
make.bat demo      # Windows

# Opens demo-report.html showing all three scenarios:
#   safe-server        — zero warnings (baseline)
#   destructive-server — delete_repo, transfer_repo flagged
#   hidden-tools-server — read_ssh_keys, send_report not in README
```

---

## Ecosystem

```
[Before use — manual]
mcp-inspect (this tool)
  → Visualize tool surface, spot surprises
  → Human decides what to allow → configure in MCP client

[At launch — automatic]
mcp-launcher
  → Secret injection, token rotation
  → Auto-block CAT-1/CAT-2 attacks via tools/list proxy
  → Audit log
```

`mcp-inspect` is intentionally standalone — no shared code with `mcp-launcher`.

---

## License

MIT
