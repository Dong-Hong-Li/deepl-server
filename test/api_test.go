// 接口测试：直接请求本地运行的 HTTP 服务。
//
// 使用前请先启动服务：
//
//	MCP_TRANSPORT=http go run ./cmd
//
// 运行测试：
//
//	go test ./test -v
//
// 可选环境变量：
//   - API_BASE_URL：服务地址，默认 http://localhost:8080
//   - MCP_AUTH_TOKEN：若服务开启了鉴权，自动附加 Bearer Token
package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
)

type languageItem struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// baseURL 返回测试目标服务地址。
func baseURL() string {
	if u := os.Getenv("API_BASE_URL"); u != "" {
		return u
	}
	return "http://localhost:8080"
}

// doRequest 发起 HTTP 请求；若配置了 MCP_AUTH_TOKEN 则自动带上鉴权头。
func doRequest(t *testing.T, method, path string, body io.Reader) *http.Response {
	t.Helper()

	req, err := http.NewRequest(method, baseURL()+path, body)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token := os.Getenv("MCP_AUTH_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v (请先启动服务: MCP_TRANSPORT=http go run ./cmd)", method, path, err)
	}
	return resp
}

// printResponse 打印请求方法与响应结果（需 go test -v 才可见）。
func printResponse(t *testing.T, method, path string, status int, body []byte) {
	t.Helper()

	t.Logf("\n========== %s %s%s ==========", method, baseURL(), path)
	t.Logf("status: %d", status)

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, body, "", "  "); err == nil {
		t.Logf("body:\n%s", pretty.String())
	} else {
		t.Logf("body:\n%s", string(body))
	}
}

// readBody 读取响应体并打印。
func readBody(t *testing.T, method, path string, resp *http.Response) []byte {
	t.Helper()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	printResponse(t, method, path, resp.StatusCode, data)
	return data
}

// TestHealth 测试 GET /health 健康检查接口。
func TestHealth(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/health", nil)
	defer resp.Body.Close()

	data := readBody(t, http.MethodGet, "/health", resp)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if string(data) != `{"status":"ok"}` {
		t.Fatalf("body = %q", data)
	}
}

// TestGetSourceLanguages 测试 GET /api/languages/source 源语言列表接口。
func TestGetSourceLanguages(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/api/languages/source", nil)
	defer resp.Body.Close()

	data := readBody(t, http.MethodGet, "/api/languages/source", resp)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var result struct {
		Languages []struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"languages"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(result.Languages) == 0 {
		t.Fatal("languages is empty")
	}
}

// TestGetTargetLanguages 测试 GET /api/languages/target 目标语言列表接口。
func TestGetTargetLanguages(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/api/languages/target", nil)
	defer resp.Body.Close()

	data := readBody(t, http.MethodGet, "/api/languages/target", resp)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var result struct {
		Languages []struct {
			Code string `json:"code"`
		} `json:"languages"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(result.Languages) == 0 {
		t.Fatal("languages is empty")
	}
}

// TestTranslateText 测试 POST /api/translate 文本翻译接口。
func TestTranslateText(t *testing.T) {
	// 英语 → 德语
	body := bytes.NewBufferString(`{"text":"hello","targetLangCode":"de"}`)
	resp := doRequest(t, http.MethodPost, "/api/translate", body)
	defer resp.Body.Close()

	data := readBody(t, http.MethodPost, "/api/translate", resp)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var result struct {
		Text               string `json:"text"`
		DetectedSourceLang string `json:"detectedSourceLang"`
		TargetLangCode     string `json:"targetLangCode"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if result.Text == "" {
		t.Fatal("text is empty")
	}
}

func codeset(langs []languageItem) map[string]bool {
	set := make(map[string]bool, len(langs))
	for _, lang := range langs {
		set[lang.Code] = true
	}
	return set
}
