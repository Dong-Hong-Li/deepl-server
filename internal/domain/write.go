package domain

// WritingStyles DeepL Write API 支持的写作风格（与官方 deepl-mcp-server 一致）。
var WritingStyles = []string{
	"academic",
	"business",
	"casual",
	"simple",
	"default",
	"prefer_academic",
	"prefer_business",
	"prefer_casual",
	"prefer_simple",
}

// WritingTones DeepL Write API 支持的语气（与官方 deepl-mcp-server 一致）。
var WritingTones = []string{
	"confident",
	"diplomatic",
	"enthusiastic",
	"friendly",
	"default",
	"prefer_confident",
	"prefer_diplomatic",
	"prefer_enthusiastic",
	"prefer_friendly",
}
