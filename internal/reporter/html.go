package reporter

import (
	"html/template"
	"io"
	"strings"
	"time"

	"github.com/masuda-masuo/mcp-inspect/internal/classifier"
	"github.com/masuda-masuo/mcp-inspect/internal/i18n"
)

// buildTemplateFuncs creates template.FuncMap bound to the given i18n.Strings.
func buildTemplateFuncs(t i18n.Strings) template.FuncMap {
	return template.FuncMap{
		"T": func(key string) string {
			switch key {
			case "TitleSuffix":
				return t.TitleSuffix
			case "LabelServers":
				return t.LabelServers
			case "LabelTools":
				return t.LabelTools
			case "LabelFlagged":
				return t.LabelFlagged
			case "NoLaunchBannerTitle":
				return t.NoLaunchBannerTitle
			case "NoLaunchBannerBody":
				return t.NoLaunchBannerBody
			case "NoticeBannerTitle":
				return t.NoticeBannerTitle
			case "NoticeBannerBody":
				return t.NoticeBannerBody
			case "ConfigPanelToggle":
				return t.ConfigPanelToggle
			case "SectionCommand":
				return t.SectionCommand
			case "SectionAllowedDirs":
				return t.SectionAllowedDirs
			case "SectionEnvVars":
				return t.SectionEnvVars
			case "ColToolName":
				return t.ColToolName
			case "ColDescription":
				return t.ColDescription
			case "ColBadges":
				return t.ColBadges
			case "EmptyTools":
				return t.EmptyTools
			case "ErrorPrefix":
				return t.ErrorPrefix
			case "NoLaunchLabel":
				return t.NoLaunchLabel
			case "ToolListToggle":
				return t.ToolListToggle
			}
			return key
		},
		"warnLabel": func(w classifier.Warning) string {
			switch w {
			case classifier.WarnDestructive:
				return t.BadgeDestructive
			case classifier.WarnWrite:
				return t.BadgeWrite
			case classifier.WarnExternalSend:
				return t.BadgeSend
			}
			return string(w)
		},
		"warnClass": func(w classifier.Warning) string {
			switch w {
			case classifier.WarnDestructive:
				return "badge-destructive"
			case classifier.WarnWrite:
				return "badge-write"
			case classifier.WarnExternalSend:
				return "badge-send"
			}
			return "badge-unknown"
		},
		"hasError": func(s ServerReport) bool { return s.Error != "" },
		"errorHint": func(s ServerReport) string {
			e := s.Error
			switch {
			case strings.Contains(e, "-32601"):
				return t.HintMethod
			case strings.Contains(e, "-32602"):
				return t.HintParams
			case strings.Contains(e, "-32600"):
				return t.HintRequest
			case strings.Contains(e, "timeout"):
				return t.HintTimeout
			case strings.Contains(e, "start:"):
				return t.HintStart
			}
			return ""
		},
		"joinStrings":  func(ss []string) string { return strings.Join(ss, " ") },
		"runtimeIcon": func(kind string) string {
			switch kind {
			case "npx", "node":
				return "\u2b61"
			case "uvx", "python":
				return "\U0001f40d"
			case "binary":
				return "\u2699"
			}
			return "\u25b6"
		},
		"runtimeLabel": func(kind string) string { return kind },
		"now":          func() string { return time.Now().Format("2006-01-02 15:04:05 UTC") },
		"breakdownLabel": func(w classifier.Warning) string {
			switch w {
			case classifier.WarnDestructive:
				return t.BadgeDestructive
			case classifier.WarnWrite:
				return t.BadgeWrite
			case classifier.WarnExternalSend:
				return t.BadgeSend
			}
			return string(w)
		},
	}
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="{{.Lang}}">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>MCP Inspect Report</title>
<style>
  :root {
    --bg: #f8f9fa; --card: #ffffff; --border: #dee2e6;
    --text: #212529; --text-muted: #6c757d; --primary: #0d6efd;
    --destructive-bg: #fff3cd; --destructive-border: #ffc107; --destructive-text: #664d03;
    --write-bg: #d1ecf1; --write-border: #0dcaf0; --write-text: #055160;
    --send-bg: #e2d9f3; --send-border: #6f42c1; --send-text: #3d1a78;
    --error-bg: #f8d7da; --error-border: #f5c2c7; --error-text: #842029;
    --nl-bg: #f0f4ff; --nl-border: #c7d2fe; --nl-text: #3730a3;
    --radius: 8px;
  }
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body { background: var(--bg); color: var(--text); font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; font-size: 14px; line-height: 1.6; }
  .container { max-width: 960px; margin: 0 auto; padding: 24px 16px; }
  .header { margin-bottom: 24px; padding-bottom: 16px; border-bottom: 1px solid var(--border); }
  .header h1 { font-size: 22px; font-weight: 700; }
  .header-meta { margin-top: 4px; color: var(--text-muted); font-size: 13px; }
  .summary { display: flex; gap: 16px; flex-wrap: wrap; margin-bottom: 24px; }
  .summary-card { background: var(--card); border: 1px solid var(--border); border-radius: var(--radius); padding: 12px 20px; flex: 1; min-width: 120px; }
  .summary-card .num { font-size: 28px; font-weight: 700; color: var(--primary); }
  .summary-card .label { font-size: 12px; color: var(--text-muted); margin-top: 2px; }
  .notice { background: #fff8e1; border: 1px solid #ffe082; border-radius: var(--radius); padding: 10px 16px; margin-bottom: 24px; font-size: 13px; color: #5d4037; }
  .notice strong { font-weight: 600; }
  .nolaunch-banner { background: var(--nl-bg); border: 1px solid var(--nl-border); border-radius: var(--radius); padding: 10px 16px; margin-bottom: 24px; font-size: 13px; color: var(--nl-text); }
  .server { background: var(--card); border: 1px solid var(--border); border-radius: var(--radius); margin-bottom: 20px; overflow: hidden; }
  .server-header { padding: 12px 16px; border-bottom: 1px solid var(--border); display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
  .server-name { font-weight: 700; font-size: 15px; }
  .server-command { font-size: 12px; color: var(--text-muted); font-family: monospace; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 480px; }
  .server-meta { margin-left: auto; font-size: 12px; color: var(--text-muted); white-space: nowrap; }
  .rt-badge { display: inline-flex; align-items: center; gap: 4px; padding: 1px 7px; border-radius: 10px; font-size: 11px; font-weight: 600; border: 1px solid; }
  .rt-npx, .rt-node { background: #dcfce7; border-color: #16a34a; color: #14532d; }
  .rt-uvx, .rt-python { background: #ede9fe; border-color: #7c3aed; color: #4c1d95; }
  .rt-binary { background: #f1f5f9; border-color: #94a3b8; color: #334155; }
  .server-error { padding: 12px 16px; background: var(--error-bg); border-left: 4px solid var(--error-border); color: var(--error-text); font-size: 13px; font-family: monospace; }
  .error-hint { font-family: -apple-system, sans-serif; margin-top: 6px; font-size: 12px; color: #5d4037; opacity: 0.85; }
  .nl-panel { padding: 0; }
  .nl-section { padding: 12px 16px; border-bottom: 1px solid var(--border); }
  .nl-section:last-child { border-bottom: none; }
  .nl-section-title { font-size: 11px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; color: var(--text-muted); margin-bottom: 8px; }
  .nl-pkg { font-family: monospace; font-size: 14px; font-weight: 600; color: var(--text); }
  .nl-cmd { font-family: monospace; font-size: 12px; color: var(--text-muted); margin-top: 2px; word-break: break-all; }
  .nl-env-key { display: inline-block; background: #e0e7ff; border-radius: 3px; padding: 1px 6px; font-size: 11px; font-family: monospace; color: var(--nl-text); margin: 2px 4px 2px 0; }
  .path-list { list-style: none; margin: 0; padding: 0; }
  .path-list li { font-family: monospace; font-size: 12px; color: var(--text); padding: 3px 0; border-bottom: 1px solid var(--border); display: flex; align-items: baseline; gap: 6px; }
  .path-list li:last-child { border-bottom: none; }
  .path-list li::before { content: "\U0001f4c1"; font-size: 11px; flex-shrink: 0; }
  .config-details { border-top: 1px solid var(--border); }
  .config-details summary { padding: 8px 16px; font-size: 12px; color: var(--text-muted); cursor: pointer; list-style: none; display: flex; align-items: center; gap: 6px; user-select: none; }
  .config-details summary::-webkit-details-marker { display: none; }
  .config-details summary::before { content: "\u25b6"; font-size: 9px; transition: transform 0.15s; }
  .config-details[open] summary::before { transform: rotate(90deg); }
  .config-details summary:hover { background: #f8f9fa; color: var(--text); }
  .config-detail-body { border-top: 1px solid var(--border); }
  table { width: 100%; border-collapse: collapse; }
  th { text-align: left; padding: 8px 16px; font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em; color: var(--text-muted); border-bottom: 1px solid var(--border); }
  td { padding: 10px 16px; border-bottom: 1px solid var(--border); vertical-align: top; }
  tr:last-child td { border-bottom: none; }
  tr:hover td { background: #f8f9fa; }
  .tool-name { font-family: monospace; font-size: 13px; font-weight: 600; }
  .tool-desc { color: var(--text-muted); font-size: 13px; }
  .tool-warnings { white-space: nowrap; }
  .empty { padding: 20px 16px; color: var(--text-muted); font-size: 13px; font-style: italic; }
  .badge { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 11px; font-weight: 600; border: 1px solid; margin-right: 4px; }
  .badge-destructive { background: var(--destructive-bg); border-color: var(--destructive-border); color: var(--destructive-text); }
  .badge-write { background: var(--write-bg); border-color: var(--write-border); color: var(--write-text); }
  .badge-send { background: var(--send-bg); border-color: var(--send-border); color: var(--send-text); }
  .footer { text-align: center; font-size: 12px; color: var(--text-muted); margin-top: 32px; }
  .tool-details summary { padding: 10px 16px; font-size: 12px; color: var(--text-muted); cursor: pointer; list-style: none; display: flex; align-items: center; gap: 8px; user-select: none; border-top: 1px solid var(--border); }
  .tool-details summary::-webkit-details-marker { display: none; }
  .tool-details summary::before { content: "\u25b6"; font-size: 9px; transition: transform 0.15s; }
  .tool-details[open] summary::before { transform: rotate(90deg); }
  .tool-details summary:hover { background: #f8f9fa; color: var(--text); }
  .tool-details .warn-count { display: inline-block; background: var(--destructive-bg); border: 1px solid var(--destructive-border); color: var(--destructive-text); border-radius: 10px; padding: 0 7px; font-size: 11px; font-weight: 600; }
</style>
</head>
<body>
<div class="container">
  <div class="header">
    <h1>\U0001f50d MCP Inspect Report{{if .NoLaunch}} <span style="font-size:14px;font-weight:400;color:#6c757d;">{{T "TitleSuffix"}}</span>{{end}}</h1>
    <div class="header-meta">{{.GeneratedAt.Format "2006-01-02 15:04:05"}} UTC &nbsp;|&nbsp; <code>{{.ConfigPath}}</code></div>
  </div>
  <div class="summary">
    <div class="summary-card"><div class="num">{{len .Servers}}</div><div class="label">{{T "LabelServers"}}</div></div>
    {{if not .NoLaunch}}
    <div class="summary-card"><div class="num">{{.TotalTools}}</div><div class="label">{{T "LabelTools"}}</div></div>
    <div class="summary-card"><div class="num">{{.TotalWarnings}}</div><div class="label">{{T "LabelFlagged"}}</div></div>
    {{range $w, $n := .WarnBreakdown}}
    <div class="summary-card"><div class="num">{{$n}}</div><div class="label">{{breakdownLabel $w}}</div></div>
    {{end}}
    {{end}}
  </div>
  {{if .NoLaunch}}
  <div class="nolaunch-banner">\U0001f4cb <strong>{{T "NoLaunchBannerTitle"}}</strong> {{T "NoLaunchBannerBody"}}</div>
  {{else}}
  <div class="notice"><strong>{{T "NoticeBannerTitle"}}</strong> {{T "NoticeBannerBody"}}</div>
  {{end}}
  {{range .Servers}}
  <div class="server">
    <div class="server-header">
      <span class="server-name">{{.Name}}</span>
      {{if .NoLaunch}}
        <span class="rt-badge rt-{{.RuntimeKind}}">{{runtimeIcon .RuntimeKind}} {{runtimeLabel .RuntimeKind}}</span>
        {{if .PackageName}}<span class="server-command">{{.PackageName}}</span>{{end}}
      {{else}}
        <span class="server-command">{{.Command}}{{if .Args}} {{joinStrings .Args}}{{end}}</span>
        <span class="server-meta">{{.ToolCount}} tools{{if gt .WarnCount 0}} \u00b7 \u26a0 {{.WarnCount}}{{end}}</span>
      {{end}}
    </div>
    {{if .NoLaunch}}
    <div class="nl-panel">
      <div class="nl-section">
        <div class="nl-section-title">{{T "SectionCommand"}}</div>
        {{if .PackageName}}<div class="nl-pkg">{{.PackageName}}</div>{{end}}
        <div class="nl-cmd">{{.Command}}{{if .Args}} {{joinStrings .Args}}{{end}}</div>
      </div>
      {{if .AllowedPaths}}
      <div class="nl-section">
        <div class="nl-section-title">{{T "SectionAllowedDirs"}} ({{len .AllowedPaths}})</div>
        <ul class="path-list">{{range .AllowedPaths}}<li>{{.}}</li>{{end}}</ul>
      </div>
      {{end}}
      {{if .EnvKeys}}
      <div class="nl-section">
        <div class="nl-section-title">{{T "SectionEnvVars"}} ({{len .EnvKeys}})</div>
        {{range .EnvKeys}}<span class="nl-env-key">{{.}}</span>{{end}}
      </div>
      {{end}}
    </div>
    {{else}}
    <details class="config-details">
      <summary>
        <span class="rt-badge rt-{{.RuntimeKind}}">{{runtimeIcon .RuntimeKind}} {{runtimeLabel .RuntimeKind}}</span>
        {{if .PackageName}}&nbsp;{{.PackageName}}{{end}}&nbsp;\u2014 {{T "ConfigPanelToggle"}}
      </summary>
      <div class="config-detail-body">
        <div class="nl-panel">
          <div class="nl-section">
            <div class="nl-section-title">{{T "SectionCommand"}}</div>
            {{if .PackageName}}<div class="nl-pkg">{{.PackageName}}</div>{{end}}
            <div class="nl-cmd">{{.Command}}{{if .Args}} {{joinStrings .Args}}{{end}}</div>
          </div>
          {{if .AllowedPaths}}
          <div class="nl-section">
            <div class="nl-section-title">{{T "SectionAllowedDirs"}} ({{len .AllowedPaths}})</div>
            <ul class="path-list">{{range .AllowedPaths}}<li>{{.}}</li>{{end}}</ul>
          </div>
          {{end}}
          {{if .EnvKeys}}
          <div class="nl-section">
            <div class="nl-section-title">{{T "SectionEnvVars"}} ({{len .EnvKeys}})</div>
            {{range .EnvKeys}}<span class="nl-env-key">{{.}}</span>{{end}}
          </div>
          {{end}}
        </div>
      </div>
    </details>
    {{if hasError .}}
    <div class="server-error">
      <div>{{T "ErrorPrefix"}} {{.Error}}</div>
      {{with errorHint .}}<div class="error-hint">\U0001f4a1 {{.}}</div>{{end}}
    </div>
    {{else if eq (len .Tools) 0}}
    <div class="empty">{{T "EmptyTools"}}</div>
    {{else}}
    <details class="tool-details"{{if gt .WarnCount 0}} open{{end}}>
      <summary>
        {{T "ToolListToggle"}} ({{.ToolCount}})
        {{if gt .WarnCount 0}}<span class="warn-count">\u26a0 {{.WarnCount}}</span>{{end}}
      </summary>
      <table>
        <thead><tr>
          <th style="width:200px">{{T "ColToolName"}}</th>
          <th>{{T "ColDescription"}}</th>
          <th style="width:160px">{{T "ColBadges"}}</th>
        </tr></thead>
        <tbody>
          {{range .Tools}}
          <tr>
            <td class="tool-name">{{.Name}}</td>
            <td class="tool-desc">{{.Description}}</td>
            <td class="tool-warnings">{{range .Warnings}}<span class="badge {{warnClass .}}">{{warnLabel .}}</span>{{end}}</td>
          </tr>
          {{end}}
        </tbody>
      </table>
    </details>
    {{end}}
    {{end}}
  </div>
  {{end}}
  <div class="footer">Generated by <a href="https://github.com/masuda-masuo/mcp-inspect">mcp-inspect</a></div>
</div>
</body>
</html>`

// WriteHTML renders the report as HTML to w.
func WriteHTML(w io.Writer, r *Report) error {
	tr := i18n.Get(r.Lang)
	tmpl, err := template.New("report").Funcs(buildTemplateFuncs(tr)).Parse(htmlTemplate)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, r)
}
