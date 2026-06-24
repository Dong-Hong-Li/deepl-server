// Package binding 提供 MCP Tool 泛型包装（Handle / Exec）。
// Handle：校验入参 + 执行业务函数 + 将响应转为 CallToolResult。
// Exec：无入参 Tool，流程同上。
package binding

import (
	"context"
	"reflect"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Handle 将 func(ctx, Req) (*Resp, error) 包装为 MCP ToolHandlerFor。
func Handle[Req, Resp any](fn func(context.Context, Req) (*Resp, error)) func(context.Context, *mcp.CallToolRequest, Req) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, req Req) (*mcp.CallToolResult, any, error) {
		if err := Validate(req); err != nil {
			return nil, nil, err
		}
		resp, err := fn(ctx, req)
		if err != nil {
			return nil, nil, err
		}
		result, err := toToolResult(resp)
		if err != nil {
			return nil, nil, err
		}
		return result, nil, nil
	}
}

// Exec 将 func(ctx) (*Resp, error) 包装为无入参 MCP Tool。
func Exec[Resp any](fn func(context.Context) (*Resp, error)) func(context.Context, *mcp.CallToolRequest, struct{}) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		resp, err := fn(ctx)
		if err != nil {
			return nil, nil, err
		}
		result, err := toToolResult(resp)
		if err != nil {
			return nil, nil, err
		}
		return result, nil, nil
	}
}

// Validate 校验 struct 上 validate:"required" 标签。
func Validate(req any) error {
	v := reflect.ValueOf(req)
	t := reflect.TypeOf(req)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
		t = t.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("validate")
		if !strings.Contains(tag, "required") {
			continue
		}
		if v.Field(i).IsZero() {
			name := field.Name
			if j := field.Tag.Get("json"); j != "" {
				name = strings.Split(j, ",")[0]
			}
			return &ValidationError{Field: name}
		}
	}
	return nil
}

// ValidationError 表示参数校验失败。
type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string {
	return e.Field + " is required"
}
