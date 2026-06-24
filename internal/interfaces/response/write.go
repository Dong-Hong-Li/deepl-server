package response

import (
	"deep-map-server/internal/delivery/binding"
	"deep-map-server/internal/domain"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type WritingStylesResponse struct {
	Styles []string
}

func NewWritingStylesResponse() *WritingStylesResponse {
	return &WritingStylesResponse{Styles: domain.WritingStyles}
}

func (r *WritingStylesResponse) MCPResult() *mcp.CallToolResult {
	return binding.TextLines(r.Styles).MCPResult()
}

type WritingTonesResponse struct {
	Tones []string
}

func NewWritingTonesResponse() *WritingTonesResponse {
	return &WritingTonesResponse{Tones: domain.WritingTones}
}

func (r *WritingTonesResponse) MCPResult() *mcp.CallToolResult {
	return binding.TextLines(r.Tones).MCPResult()
}

type RephraseTextResponse struct {
	Text string `json:"text"`
}

func (r *RephraseTextResponse) MCPResult() *mcp.CallToolResult {
	return binding.TextLines{r.Text}.MCPResult()
}

type TranslateDocumentResponse struct {
	Status           string `json:"status"`
	BilledCharacters int    `json:"billedCharacters"`
	OutputFilePath   string `json:"outputFilePath"`
	TargetLangCode   string `json:"targetLangCode"`
}

func (r *TranslateDocumentResponse) MCPResult() *mcp.CallToolResult {
	return binding.TextLines{
		"Document translated successfully! Status: " + r.Status,
		"Target language used: " + r.TargetLangCode,
		"Characters billed: " + strconv.Itoa(r.BilledCharacters),
		"Output file path: " + r.OutputFilePath,
	}.MCPResult()
}
