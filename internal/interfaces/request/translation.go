package request

// TranslateTextRequest translate-text 工具入参。
type TranslateTextRequest struct {
	// 要翻译的文本
	Text string `json:"text" jsonschema:"Text to translate" validate:"required"`
	// 源语言代码，留空则自动检测（短文本可能误判，建议显式指定）
	SourceLangCode string `json:"sourceLangCode,omitempty" jsonschema:"source language code from /api/languages/source; recommended for short text"`
	// 目标语言代码
	TargetLangCode string `json:"targetLangCode" jsonschema:"target language code" validate:"required"`
	// 语气控制：less, more, default, prefer_less, prefer_more
	Formality string `json:"formality,omitempty" jsonschema:"Controls formality: less, more, default, prefer_less, prefer_more"`
	// 术语表ID
	GlossaryID string `json:"glossaryId,omitempty" jsonschema:"Glossary ID to ensure consistent terminology translation"`
}
