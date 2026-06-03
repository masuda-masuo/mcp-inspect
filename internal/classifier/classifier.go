// Package classifier assigns warning badges to MCP tools based on
// keyword matching against tool names and descriptions.
// Accuracy is intentionally traded for simplicity: the goal is awareness,
// not a security gate.
package classifier

import (
	"strings"
)

// Warning represents a single warning category.
type Warning string

const (
	WarnDestructive  Warning = "destructive"
	WarnWrite        Warning = "write"
	WarnExternalSend Warning = "external_send"
)

// Label returns the English label (used for JSON output and as a default).
func (w Warning) Label() string {
	switch w {
	case WarnDestructive:
		return "Destructive"
	case WarnWrite:
		return "Write"
	case WarnExternalSend:
		return "External send"
	}
	return string(w)
}

var destructiveNameKeywords = []string{
	"delete", "remove", "destroy", "drop", "purge", "transfer",
	"truncate", "wipe", "erase", "revoke", "ban", "terminate",
}

var writeNameKeywords = []string{
	"write", "update", "push", "insert", "overwrite", "create",
	"edit", "modify", "patch", "upload", "save",
	"move", "rename", "copy", "merge", "replace",
	"add_", "_add", "set_", "_set", "put_", "_put",
	"commit_",
}

var externalSendNameKeywords = []string{
	"send", "publish", "notify", "webhook", "submit",
	"email", "sms", "tweet", "broadcast",
}

var destructiveDescKeywords = []string{
	"cannot be undone", "irreversible", "permanently delete",
	"取り消せません", "元に戻せません",
}

var writeDescKeywords = []string{
	"writes to", "write to", "overwrites", "modifies the file",
}

var externalSendDescKeywords = []string{
	"sends to external", "external endpoint", "外部エンドポイント",
	"外部サービスに送信",
}

// Classify returns warnings for a tool given its name and description.
func Classify(name, description string) []Warning {
	lname := strings.ToLower(name)
	ldesc := strings.ToLower(description)

	var warnings []Warning

	destructive := matchesAny(lname, destructiveNameKeywords) ||
		matchesAny(ldesc, destructiveDescKeywords)
	write := matchesAny(lname, writeNameKeywords) ||
		matchesAny(ldesc, writeDescKeywords)

	if destructive {
		warnings = append(warnings, WarnDestructive)
	} else if write {
		warnings = append(warnings, WarnWrite)
	}

	externalSend := matchesAny(lname, externalSendNameKeywords) ||
		matchesAny(ldesc, externalSendDescKeywords)
	if externalSend {
		warnings = append(warnings, WarnExternalSend)
	}

	return warnings
}

func matchesAny(s string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(s, kw) {
			return true
		}
	}
	return false
}
