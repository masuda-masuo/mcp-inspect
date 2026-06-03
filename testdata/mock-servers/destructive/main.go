// destructive-mock-server: MCP mock containing write/destructive tools.
// Demonstrates that even documented tools can be easy to overlook.
// No network or file I/O.
package main

import (
	"bufio"
	"encoding/json"
	"os"
)

var tools = []map[string]any{
	{
		"name":        "list_issues",
		"description": "GitHubリポジトリのIssue一覧を取得します。",
		"inputSchema": map[string]any{"type": "object", "properties": map[string]any{
			"repo": map[string]any{"type": "string"},
		}},
	},
	{
		"name":        "create_issue",
		"description": "GitHubにIssueを作成します。",
		"inputSchema": map[string]any{"type": "object", "properties": map[string]any{
			"repo":  map[string]any{"type": "string"},
			"title": map[string]any{"type": "string"},
			"body":  map[string]any{"type": "string"},
		}},
	},
	{
		"name":        "push_files",
		"description": "リポジトリにファイルをプッシュします。",
		"inputSchema": map[string]any{"type": "object", "properties": map[string]any{
			"repo":  map[string]any{"type": "string"},
			"files": map[string]any{"type": "array"},
		}},
	},
	{
		"name":        "delete_repo",
		"description": "GitHubリポジトリを完全に削除します。この操作は取り消せません。",
		"inputSchema": map[string]any{"type": "object", "properties": map[string]any{
			"repo": map[string]any{"type": "string"},
		}},
	},
	{
		"name":        "transfer_repo",
		"description": "リポジトリの所有権を別のアカウントまたは組織に移譲します。",
		"inputSchema": map[string]any{"type": "object", "properties": map[string]any{
			"repo":        map[string]any{"type": "string"},
			"new_owner":   map[string]any{"type": "string"},
		}},
	},
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 64*1024), 64*1024)
	for scanner.Scan() {
		handle(scanner.Text())
	}
}

func handle(line string) {
	var req map[string]any
	if err := json.Unmarshal([]byte(line), &req); err != nil {
		return
	}
	method, _ := req["method"].(string)
	id := req["id"]
	switch method {
	case "initialize":
		respond(id, map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]any{"name": "destructive-mock", "version": "0.1.0"},
		})
	case "notifications/initialized":
	case "tools/list":
		respond(id, map[string]any{"tools": tools})
	}
}

func respond(id any, result any) {
	resp := map[string]any{"jsonrpc": "2.0", "id": id, "result": result}
	data, _ := json.Marshal(resp)
	os.Stdout.Write(data)
	os.Stdout.Write([]byte("\n"))
}
