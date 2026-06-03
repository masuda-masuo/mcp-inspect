package reporter

import (
	"encoding/json"
	"io"
)

// JSONOutput is the schema for --format json.
type JSONOutput struct {
	GeneratedAt string       `json:"generated_at"`
	Config      string       `json:"config"`
	Servers     []jsonServer `json:"servers"`
}

type jsonServer struct {
	Name      string     `json:"name"`
	Command   string     `json:"command"`
	ToolCount int        `json:"tool_count"`
	Tools     []jsonTool `json:"tools"`
	Error     string     `json:"error,omitempty"`
}

type jsonTool struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Warnings    []string `json:"warnings"`
}

// WriteJSON serialises the report as JSON to w.
func WriteJSON(w io.Writer, r *Report) error {
	out := JSONOutput{
		GeneratedAt: r.GeneratedAt.Format("2006-01-02T15:04:05Z"),
		Config:      r.ConfigPath,
	}
	for _, sr := range r.Servers {
		js := jsonServer{
			Name:      sr.Name,
			Command:   sr.Command,
			ToolCount: sr.ToolCount,
			Error:     sr.Error,
		}
		for _, t := range sr.Tools {
			jt := jsonTool{
				Name:        t.Name,
				Description: t.Description,
			}
			for _, w := range t.Warnings {
				jt.Warnings = append(jt.Warnings, string(w))
			}
			if jt.Warnings == nil {
				jt.Warnings = []string{}
			}
			js.Tools = append(js.Tools, jt)
		}
		if js.Tools == nil {
			js.Tools = []jsonTool{}
		}
		out.Servers = append(out.Servers, js)
	}
	if out.Servers == nil {
		out.Servers = []jsonServer{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
