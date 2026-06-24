package binding

import (
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Responder 将业务响应转换为 MCP Tool 结果。
type Responder interface {
	MCPResult() *mcp.CallToolResult
}

// TextLines 以多段文本块返回（与官方 deepl-mcp-server 行为一致）。
type TextLines []string

func (lines TextLines) MCPResult() *mcp.CallToolResult {
	content := make([]mcp.Content, len(lines))
	for i, line := range lines {
		content[i] = &mcp.TextContent{Text: line}
	}
	return &mcp.CallToolResult{Content: content}
}

func toToolResult(resp any) (*mcp.CallToolResult, error) {
	if resp == nil {
		return &mcp.CallToolResult{}, nil
	}
	if r, ok := resp.(Responder); ok {
		return r.MCPResult(), nil
	}

	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return nil, err
	}
	return TextLines{string(data)}.MCPResult(), nil
}
