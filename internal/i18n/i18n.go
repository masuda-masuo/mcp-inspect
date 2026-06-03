// Package i18n provides UI string translations for HTML report output.
// Supported languages: "en" (default), "ja".
package i18n

// Lang represents a supported UI language.
type Lang string

const (
	EN Lang = "en"
	JA Lang = "ja"
)

// Parse normalises a language tag string into a Lang.
// Unknown values fall back to EN.
func Parse(s string) Lang {
	switch s {
	case "ja", "ja-JP":
		return JA
	default:
		return EN
	}
}

// Strings holds all translatable UI strings used in the HTML report.
type Strings struct {
	// Header
	TitleSuffix string // "(config only · no launch)"

	// Summary cards
	LabelServers string
	LabelTools   string
	LabelFlagged string

	// Banners
	NoLaunchBannerTitle string
	NoLaunchBannerBody  string
	NoticeBannerTitle   string
	NoticeBannerBody    string

	// Config panel
	ConfigPanelToggle  string
	SectionCommand     string
	SectionAllowedDirs string // prefix; count is appended
	SectionEnvVars     string // prefix; count is appended

	// Tool table
	ColToolName    string
	ColDescription string
	ColBadges      string
	EmptyTools     string

	// Error
	ErrorPrefix string

	// Error hints
	HintMethod  string
	HintParams  string
	HintRequest string
	HintTimeout string
	HintStart   string

	// Badge labels
	BadgeDestructive string
	BadgeWrite       string
	BadgeSend        string

	// Tool list toggle
	ToolListToggle string

	// No-launch mode label
	NoLaunchLabel string
}

var translations = map[Lang]Strings{
	EN: {
		TitleSuffix: "(config only \u00b7 no launch)",

		LabelServers: "Servers",
		LabelTools:   "Total tools",
		LabelFlagged: "Flagged",

		NoLaunchBannerTitle: "--no-launch mode:",
		NoLaunchBannerBody:  "Servers were not started. Tool lists are unavailable. Remove --no-launch to fetch them.",
		NoticeBannerTitle:   "\u2139\ufe0f How to use:",
		NoticeBannerBody:    "Compare this report with each server's README. Badges are awareness prompts, not blocks. Decide what to allow in your MCP client's permission settings.",

		ConfigPanelToggle:  "show config",
		SectionCommand:     "Command",
		SectionAllowedDirs: "Allowed directories",
		SectionEnvVars:     "Environment variables",

		ColToolName:    "Tool name",
		ColDescription: "Description",
		ColBadges:      "Badges",
		EmptyTools:     "No tools (tools/list returned empty)",

		ErrorPrefix: "\u26a0 Launch error:",

		HintMethod:  "This server may not implement tools/list (old protocol version).",
		HintParams:  "Request parameters were rejected. A proxy server may require authentication or additional configuration.",
		HintRequest: "Invalid request error. The server may be using a different protocol version.",
		HintTimeout: "Server startup or response timed out. Check the command path and arguments.",
		HintStart:   "Failed to start the server. Check the command path and permissions.",

		BadgeDestructive: "Destructive",
		BadgeWrite:       "Write",
		BadgeSend:        "External send",

		ToolListToggle: "Tools",
		NoLaunchLabel:  "config only",
	},
	JA: {
		TitleSuffix: "\uff08\u8a2d\u5b9a\u306e\u307f\u30fb\u8d77\u52d5\u306a\u3057\uff09",

		LabelServers: "\u30b5\u30fc\u30d0\u30fc",
		LabelTools:   "\u30c4\u30fc\u30eb\u5408\u8a08",
		LabelFlagged: "\u8981\u78ba\u8a8d",

		NoLaunchBannerTitle: "--no-launch \u30e2\u30fc\u30c9\uff1a",
		NoLaunchBannerBody:  "\u30b5\u30fc\u30d0\u30fc\u3092\u8d77\u52d5\u305b\u305a\u8a2d\u5b9a\u30d5\u30a1\u30a4\u30eb\u306e\u5185\u5bb9\u306e\u307f\u3092\u8868\u793a\u3057\u3066\u3044\u307e\u3059\u3002\u30c4\u30fc\u30eb\u4e00\u89a7\u306f\u53d6\u5f97\u3055\u308c\u307e\u305b\u3093\u3002\u8d77\u52d5\u3057\u3066\u30c4\u30fc\u30eb\u4e00\u89a7\u3092\u53d6\u5f97\u3059\u308b\u306b\u306f --no-launch \u3092\u5916\u3057\u3066\u5b9f\u884c\u3057\u3066\u304f\u3060\u3055\u3044\u3002",
		NoticeBannerTitle:   "\u2139\ufe0f \u4f7f\u3044\u65b9\uff1a",
		NoticeBannerBody:    "\u3053\u306e\u30ec\u30dd\u30fc\u30c8\u3092\u30b5\u30fc\u30d0\u30fc\u306e README \u3068\u898b\u6bd4\u3079\u3066\u304f\u3060\u3055\u3044\u3002\u30d0\u30c3\u30b8\u306f\u8b66\u544a\u3067\u306f\u306a\u304f\u300c\u6c17\u3065\u304d\u300d\u306e\u63d0\u793a\u3067\u3059\u3002\u30d6\u30ed\u30c3\u30af\u306f\u3057\u307e\u305b\u3093\u3002\u4f7f\u3046\uff0f\u4f7f\u308f\u306a\u3044\u306e\u5224\u65ad\u306f\u3001MCP\u30af\u30e9\u30a4\u30a2\u30f3\u30c8\u306e\u8a31\u53ef\u8a2d\u5b9a\u3067\u884c\u3063\u3066\u304f\u3060\u3055\u3044\u3002",

		ConfigPanelToggle:  "\u8a2d\u5b9a\u3092\u8868\u793a",
		SectionCommand:     "\u30b3\u30de\u30f3\u30c9",
		SectionAllowedDirs: "\u8a31\u53ef\u30c7\u30a3\u30ec\u30af\u30c8\u30ea",
		SectionEnvVars:     "\u74b0\u5883\u5909\u6570",

		ColToolName:    "\u30c4\u30fc\u30eb\u540d",
		ColDescription: "\u8aac\u660e",
		ColBadges:      "\u30d0\u30c3\u30b8",
		EmptyTools:     "\u30c4\u30fc\u30eb\u306a\u3057\uff08tools/list \u304c\u7a7a\u3067\u3057\u305f\uff09",

		ErrorPrefix: "\u26a0 \u8d77\u52d5\u30a8\u30e9\u30fc:",

		HintMethod:  "\u3053\u306e\u30b5\u30fc\u30d0\u30fc\u306f tools/list \u30e1\u30bd\u30c3\u30c9\u306b\u5bfe\u5fdc\u3057\u3066\u3044\u306a\u3044\u53ef\u80fd\u6027\u304c\u3042\u308a\u307e\u3059\uff08\u53e4\u3044\u30d7\u30ed\u30c8\u30b3\u30eb\u30d0\u30fc\u30b8\u30e7\u30f3\uff09\u3002",
		HintParams:  "\u30ea\u30af\u30a8\u30b9\u30c8\u30d1\u30e9\u30e1\u30fc\u30bf\u304c\u62d2\u5426\u3055\u308c\u307e\u3057\u305f\u3002\u30d7\u30ed\u30ad\u30b7\u7d4c\u7531\u306e\u30b5\u30fc\u30d0\u30fc\u3067\u8a8d\u8a3c\u3084\u8a2d\u5b9a\u304c\u5fc5\u8981\u306a\u5834\u5408\u304c\u3042\u308a\u307e\u3059\u3002",
		HintRequest: "\u4e0d\u6b63\u306a\u30ea\u30af\u30a8\u30b9\u30c8\u30a8\u30e9\u30fc\u3067\u3059\u3002\u30b5\u30fc\u30d0\u30fc\u306e\u30d7\u30ed\u30c8\u30b3\u30eb\u30d0\u30fc\u30b8\u30e7\u30f3\u304c\u7570\u306a\u308b\u53ef\u80fd\u6027\u304c\u3042\u308a\u307e\u3059\u3002",
		HintTimeout: "\u30b5\u30fc\u30d0\u30fc\u306e\u8d77\u52d5\u307e\u305f\u306f\u30ec\u30b9\u30dd\u30f3\u30b9\u304c\u30bf\u30a4\u30e0\u30a2\u30a6\u30c8\u3057\u307e\u3057\u305f\u3002\u30b3\u30de\u30f3\u30c9\u30d1\u30b9\u3084\u5f15\u6570\u3092\u78ba\u8a8d\u3057\u3066\u304f\u3060\u3055\u3044\u3002",
		HintStart:   "\u30b5\u30fc\u30d0\u30fc\u306e\u8d77\u52d5\u306b\u5931\u6557\u3057\u307e\u3057\u305f\u3002\u30b3\u30de\u30f3\u30c9\u30d1\u30b9\u30fb\u5b9f\u884c\u6a29\u9650\u3092\u78ba\u8a8d\u3057\u3066\u304f\u3060\u3055\u3044\u3002",

		BadgeDestructive: "\u7834\u58ca\u7684\u64cd\u4f5c",
		BadgeWrite:       "\u66f8\u304d\u8fbc\u307f\u64cd\u4f5c",
		BadgeSend:        "\u5916\u90e8\u9001\u4fe1",

		ToolListToggle: "\u30c4\u30fc\u30eb\u4e00\u89a7",
		NoLaunchLabel:  "\u8a2d\u5b9a\u306e\u307f",
	},
}

// Get returns the Strings for the given Lang.
func Get(l Lang) Strings {
	if s, ok := translations[l]; ok {
		return s
	}
	return translations[EN]
}
