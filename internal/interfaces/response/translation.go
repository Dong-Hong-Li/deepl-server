package response

import (
	"encoding/json"

	"deep-map-server/internal/delivery/binding"
	"deep-map-server/internal/domain"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// LanguageListResponse 语言列表响应。
type LanguageListResponse struct {
	Languages []domain.Language `json:"languages"`
}

func NewLanguageList(languages []domain.Language) *LanguageListResponse {
	return &LanguageListResponse{Languages: languages}
}

func (r *LanguageListResponse) MCPResult() *mcp.CallToolResult {
	lines := make([]string, len(r.Languages))
	for i, lang := range r.Languages {
		data, _ := json.Marshal(lang)
		lines[i] = string(data)
	}
	return binding.TextLines(lines).MCPResult()
}

// TranslateTextResponse 文本翻译响应。
type TranslateTextResponse struct {
	Text               string `json:"text"`
	DetectedSourceLang string `json:"detectedSourceLang"`
	TargetLangCode     string `json:"targetLangCode"`
}

func (r *TranslateTextResponse) MCPResult() *mcp.CallToolResult {
	return binding.TextLines{
		r.Text,
		"Detected source language: " + r.DetectedSourceLang,
		"Target language used: " + r.TargetLangCode,
	}.MCPResult()
}
