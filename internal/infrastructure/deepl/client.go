package deepl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ProBaseURL  = "https://api.deepl.com"
	FreeBaseURL = "https://api-free.deepl.com"

	AuthHeader = "Authorization"
	AuthScheme = "DeepL-Auth-Key"
)

// Client 是 DeepL 底层 HTTP 客户端，只负责请求发送与响应读取。
// 具体 API 语义请放在 adapter 中实现。
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

type Option func(*Client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.http = httpClient
	}
}

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

func NewClient(apiKey string, opts ...Option) *Client {
	client := &Client{
		apiKey:  apiKey,
		baseURL: ResolveBaseURL(apiKey),
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// ResolveBaseURL 根据 API Key 判断 Free / Pro 端点。
func ResolveBaseURL(apiKey string) string {
	if strings.HasSuffix(apiKey, ":fx") {
		return FreeBaseURL
	}
	return ProBaseURL
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

// Get 发送 GET 请求并返回响应体。
func (c *Client) Get(ctx context.Context, path string, query url.Values) ([]byte, error) {
	endpoint := c.baseURL + path
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	c.setAuth(req)

	return c.doRaw(req)
}

/**
* @description: Post form data to the given path.
* @param {context.Context} ctx
* @param {string} path
* @param {url.Values} form
* @return {[]byte, error}
 */
func (c *Client) PostForm(ctx context.Context, path string, form url.Values) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+path,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, err
	}
	c.setAuth(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.doRaw(req)
}

// PostJSON 发送 JSON POST 请求并返回响应体。
func (c *Client) PostJSON(ctx context.Context, path string, body any) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	c.setAuth(req)
	req.Header.Set("Content-Type", "application/json")

	return c.doRaw(req)
}

// PostMultipart 上传 multipart/form-data 并返回响应体。
func (c *Client) PostMultipart(ctx context.Context, path string, fields map[string]string, fileField, filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	part, err := writer.CreateFormFile(fileField, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, &body)
	if err != nil {
		return nil, err
	}
	c.setAuth(req)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return c.doRaw(req)
}

func (c *Client) setAuth(req *http.Request) {
	req.Header.Set(AuthHeader, AuthScheme+" "+c.apiKey)
}

func (c *Client) doRaw(req *http.Request) ([]byte, error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       strings.TrimSpace(string(body)),
		}
	}
	return body, nil
}

// PostFormPath 发送 form POST 到指定 path。
func (c *Client) PostFormPath(ctx context.Context, path string, form url.Values) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+path,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, err
	}
	c.setAuth(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.doRaw(req)
}

// APIError 表示 DeepL HTTP 错误响应。
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("deepl API error (%d): %s", e.StatusCode, e.Body)
}
