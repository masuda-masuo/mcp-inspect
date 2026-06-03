// safe-mock-server: baseline MCP mock with zero warnings.
// Returns read-only filesystem tools only. No network or file I/O.
package main

import (
	"bufio"
	"encoding/json"
	"os"
)

var tools = []map[string]any{
	{
		"name":        "read_file",
		"description": "ファイルの内容を読み取ります。",
		"inputSchema": map[string]any{"type": "object", "properties": map[string]any{
			"path": map[string]any{"type": "string", "description": "ファイルパス"},
		}},
	},
	{
		"name":        "list_directory",
		"description": "ディレクトリ内のファイル一覧を取得します。",
		"inputSchema": map[string]any{"type": "object", "properties": map[string]any{
			"path": map[string]any{"type": "string"},
		}},
	},
	{
		"name":        "get_file_info",
		"description": "ファイルのメタデータ（サイズ、更新日時など）を取得します。",
		"inputSchema": map[string]any{"type": "object", "properties": map[string]any{
			"path": map[string]any{"type": "string"},
		}},
	},
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 64*1024), 64*1024)
	for scanner.Scan() {
		line := scanner.Text()
		handle(line)
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
			"serverInfo":      map[string]any{"name": "safe-mock", "version": "0.1.0"},
		})
	case "notifications/initialized":
		// no response required for notifications
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
