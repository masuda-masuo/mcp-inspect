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
	"\u53d6\u308a\u6d88\u305b\u307e\u305b\u3093", "\u5143\u306b\u623b\u305b\u307e\u305b\u3093",
}

var writeDescKeywords = []string{
	"writes to", "write to", "overwrites", "modifies the file",
}

var externalSendDescKeywords = []string{
	"sends to external", "external endpoint", "\u5916\u90e8\u30a8\u30f3\u30c9\u30dd\u30a4\u30f3\u30c8",
	"\u5916\u90e8\u30b5\u30fc\u30d3\u30b9\u306b\u9001\u4fe1",
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
