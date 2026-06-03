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
	LabelServers  string
	LabelTools    string
	LabelFlagged  string

	// Banners
	NoLaunchBannerTitle string
	NoLaunchBannerBody  string
	NoticeBannerTitle   string
	NoticeBannerBody    string

	// Config panel
	ConfigPanelToggle    string
	SectionCommand       string
	SectionAllowedDirs   string // prefix; count is appended
	SectionEnvVars       string // prefix; count is appended

	// Tool table
	ColToolName    string
	ColDescription string
	ColBadges      string
	EmptyTools     string

	// Error
	ErrorPrefix string

	// Error hints
	HintMethod    string
	HintParams    string
	HintRequest   string
	HintTimeout   string
	HintStart     string

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
		TitleSuffix: "(config only · no launch)",

		LabelServers: "Servers",
		LabelTools:   "Total tools",
		LabelFlagged: "Flagged",

		NoLaunchBannerTitle: "--no-launch mode:",
		NoLaunchBannerBody:  "Servers were not started. Tool lists are unavailable. Remove --no-launch to fetch them.",
		NoticeBannerTitle:   "ℹ️ How to use:",
		NoticeBannerBody:    "Compare this report with each server's README. Badges are awareness prompts, not blocks. Decide what to allow in your MCP client's permission settings.",

		ConfigPanelToggle:  "show config",
		SectionCommand:     "Command",
		SectionAllowedDirs: "Allowed directories",
		SectionEnvVars:     "Environment variables",

		ColToolName:    "Tool name",
		ColDescription: "Description",
		ColBadges:      "Badges",
		EmptyTools:     "No tools (tools/list returned empty)",

		ErrorPrefix: "⚠ Launch error:",

		HintMethod:  "This server may not implement tools/list (old protocol version).",
		HintParams:  "Request parameters were rejected. A proxy server may require authentication or additional configuration.",
		HintRequest: "Invalid request error. The server may be using a different protocol version.",
		HintTimeout: "Server startup or response timed out. Check the command path and arguments.",
		HintStart:   "Failed to start the server. Check the command path and permissions.",

		BadgeDestructive: "Destructive",
		BadgeWrite:       "Write",
		BadgeSend:        "External send",

		ToolListToggle: "Tools",
		NoLaunchLabel: "config only",
	},
	JA: {
		TitleSuffix: "（設定のみ・起動なし）",

		LabelServers: "サーバー",
		LabelTools:   "ツール合計",
		LabelFlagged: "要確認",

		NoLaunchBannerTitle: "--no-launch モード：",
		NoLaunchBannerBody:  "サーバーを起動せず設定ファイルの内容のみを表示しています。ツール一覧は取得されません。起動してツール一覧を取得するには --no-launch を外して実行してください。",
		NoticeBannerTitle:   "ℹ️ 使い方：",
		NoticeBannerBody:    "このレポートをサーバーの README と見比べてください。バッジは警告ではなく「気づき」の提示です。ブロックはしません。使う／使わないの判断は、MCPクライアントの許可設定で行ってください。",

		ConfigPanelToggle:  "設定を表示",
		SectionCommand:     "コマンド",
		SectionAllowedDirs: "許可ディレクトリ",
		SectionEnvVars:     "環境変数",

		ColToolName:    "ツール名",
		ColDescription: "説明",
		ColBadges:      "バッジ",
		EmptyTools:     "ツールなし（tools/list が空でした）",

		ErrorPrefix: "⚠ 起動エラー:",

		HintMethod:  "このサーバーは tools/list メソッドに対応していない可能性があります（古いプロトコルバージョン）。",
		HintParams:  "リクエストパラメータが拒否されました。プロキシ経由のサーバーで認証や設定が必要な場合があります。",
		HintRequest: "不正なリクエストエラーです。サーバーのプロトコルバージョンが異なる可能性があります。",
		HintTimeout: "サーバーの起動またはレスポンスがタイムアウトしました。コマンドパスや引数を確認してください。",
		HintStart:   "サーバーの起動に失敗しました。コマンドパス・実行権限を確認してください。",

		BadgeDestructive: "破壊的操作",
		BadgeWrite:       "書き込み操作",
		BadgeSend:        "外部送信",

		ToolListToggle: "ツール一覧",
		NoLaunchLabel: "設定のみ",
	},
}

// Get returns the Strings for the given Lang.
func Get(l Lang) Strings {
	if s, ok := translations[l]; ok {
		return s
	}
	return translations[EN]
}
