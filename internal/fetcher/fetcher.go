// Package fetcher launches an MCP server and retrieves its tool list via stdio JSON-RPC.
//
// Security note: to obtain tool descriptions we must start the server process.
// The server's main() runs before we can read tools/list, so side-effects
// (file access, network calls) in the server binary cannot be prevented here.
// See README for the implications.
package fetcher

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/masuda-masuo/mcp-inspect/internal/config"
)

// Tool represents a single MCP tool as returned by tools/list.
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Result holds the outcome of inspecting one server.
type Result struct {
	ServerName string
	Command    string
	Tools      []Tool
	Error      string
}

const dialTimeout = 15 * time.Second

// FetchAll starts each server in cfg and retrieves its tool list.
func FetchAll(cfg *config.Config) []Result {
	results := make([]Result, 0, len(cfg.Servers))
	for name, srv := range cfg.Servers {
		r := fetchOne(name, srv)
		results = append(results, r)
	}
	return results
}

// resolveCommand resolves the executable path.
// On Windows, if the command has no extension and no .exe variant is found
// by PATH lookup, we try appending .exe explicitly so that relative paths
// like "testdata/mock-servers/safe/mock-server" just work.
func resolveCommand(command string) string {
	if runtime.GOOS != "windows" {
		return command
	}
	lower := strings.ToLower(command)
	if strings.HasSuffix(lower, ".exe") ||
		strings.HasSuffix(lower, ".cmd") ||
		strings.HasSuffix(lower, ".bat") {
		return command
	}
	// Try PATH lookup first (handles system commands like "npx", "node", etc.)
	if _, err := exec.LookPath(command); err == nil {
		return command
	}
	// Fall back to appending .exe (handles relative paths to compiled binaries)
	return command + ".exe"
}

func fetchOne(name string, srv config.ServerConfig) Result {
	r := Result{
		ServerName: name,
		Command:    buildCommandLine(srv),
	}

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()

	resolvedCmd := resolveCommand(srv.Command)
	args := append([]string{}, srv.Args...)
	cmd := exec.CommandContext(ctx, resolvedCmd, args...)

	// Inject custom env vars on top of the current environment.
	if len(srv.Env) > 0 {
		base := os.Environ()
		for k, v := range srv.Env {
			base = append(base, k+"="+v)
		}
		cmd.Env = base
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		r.Error = fmt.Sprintf("stdin pipe: %v", err)
		return r
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		r.Error = fmt.Sprintf("stdout pipe: %v", err)
		return r
	}

	if err := cmd.Start(); err != nil {
		r.Error = fmt.Sprintf("start: %v", err)
		return r
	}
	defer func() {
		_ = stdin.Close()
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	tools, err := negotiateAndList(stdin, stdout, ctx)
	if err != nil {
		r.Error = err.Error()
		return r
	}
	r.Tools = tools
	return r
}

// ---------- JSON-RPC helpers ----------

// jsonrpcRequest is used for requests (with ID) and notifications (ID omitted via omitempty).
type jsonrpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      *int   `json:"id,omitempty"` // nil → notification (no id field in JSON)
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

type jsonrpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *jsonrpcError   `json:"error"`
}

type jsonrpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func intPtr(i int) *int { return &i }

func sendRequest(w io.Writer, req jsonrpcRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n", data)
	return err
}

// readResponse reads the next non-empty line and parses it as a JSON-RPC response.
// It skips lines that look like notifications (no "id" field or id==null) from the server.
func readResponse(scanner *bufio.Scanner, ctx context.Context) (*jsonrpcResponse, error) {
	type result struct {
		resp *jsonrpcResponse
		err  error
	}
	ch := make(chan result, 1)

	go func() {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			var resp jsonrpcResponse
			if err := json.Unmarshal([]byte(line), &resp); err != nil {
				ch <- result{nil, fmt.Errorf("invalid JSON from server: %v (line: %s)", err, line)}
				return
			}
			// Skip server-initiated notifications (id is null or absent)
			if resp.ID == nil {
				continue
			}
			ch <- result{&resp, nil}
			return
		}
		ch <- result{nil, fmt.Errorf("server closed stdout without responding")}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout waiting for response")
	case r := <-ch:
		return r.resp, r.err
	}
}

// negotiateAndList performs the MCP initialize handshake and then calls tools/list.
func negotiateAndList(stdin io.Writer, stdout io.Reader, ctx context.Context) ([]Tool, error) {
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	// 1. initialize
	initReq := jsonrpcRequest{
		JSONRPC: "2.0",
		ID:      intPtr(1),
		Method:  "initialize",
		Params: map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{},
			"clientInfo": map[string]any{
				"name":    "mcp-inspect",
				"version": "0.1.0",
			},
		},
	}
	if err := sendRequest(stdin, initReq); err != nil {
		return nil, fmt.Errorf("sending initialize: %v", err)
	}

	initResp, err := readResponse(scanner, ctx)
	if err != nil {
		return nil, fmt.Errorf("initialize response: %v", err)
	}
	if initResp.Error != nil {
		return nil, fmt.Errorf("initialize error %d: %s", initResp.Error.Code, initResp.Error.Message)
	}

	// Check if server advertises tools capability
	var initResult struct {
		Capabilities struct {
			Tools *json.RawMessage `json:"tools"`
		} `json:"capabilities"`
	}
	if err := json.Unmarshal(initResp.Result, &initResult); err == nil {
		if initResult.Capabilities.Tools == nil {
			// Server doesn't advertise tools capability – return empty gracefully
			return []Tool{}, nil
		}
	}

	// 2. initialized notification – must have NO id field (it's a notification, not a request)
	notif := jsonrpcRequest{
		JSONRPC: "2.0",
		// ID is nil (omitempty) → not included in JSON
		Method: "notifications/initialized",
	}
	if err := sendRequest(stdin, notif); err != nil {
		return nil, fmt.Errorf("sending initialized notification: %v", err)
	}

	// 3. tools/list (paginated)
	// Send params only when we have a cursor; omit params entirely on first call
	// to maximise compatibility with older servers.
	var allTools []Tool
	var cursor *string
	reqID := 2

	for {
		var params any
		if cursor != nil {
			params = map[string]any{"cursor": *cursor}
		}
		// params == nil → omitempty removes it from JSON

		listReq := jsonrpcRequest{
			JSONRPC: "2.0",
			ID:      intPtr(reqID),
			Method:  "tools/list",
			Params:  params,
		}
		reqID++

		if err := sendRequest(stdin, listReq); err != nil {
			return nil, fmt.Errorf("sending tools/list: %v", err)
		}

		listResp, err := readResponse(scanner, ctx)
		if err != nil {
			return nil, fmt.Errorf("tools/list response: %v", err)
		}
		if listResp.Error != nil {
			return nil, fmt.Errorf("tools/list error %d: %s", listResp.Error.Code, listResp.Error.Message)
		}

		var result struct {
			Tools      []Tool  `json:"tools"`
			NextCursor *string `json:"nextCursor"`
		}
		if err := json.Unmarshal(listResp.Result, &result); err != nil {
			return nil, fmt.Errorf("decoding tools/list result: %v", err)
		}

		allTools = append(allTools, result.Tools...)
		if result.NextCursor == nil || *result.NextCursor == "" {
			break
		}
		cursor = result.NextCursor
	}

	return allTools, nil
}

func buildCommandLine(srv config.ServerConfig) string {
	parts := []string{srv.Command}
	parts = append(parts, srv.Args...)
	return strings.Join(parts, " ")
}
