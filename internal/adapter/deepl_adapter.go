package adapter

import (
	"context"
	"deep-map-server/internal/domain"
	"deep-map-server/internal/infrastructure/deepl"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

const languagesPath = "/v2/languages"
const translatePath = "/v2/translate"
const rephrasePath = "/v2/write/rephrase"
const documentPath = "/v2/document"

const sourceLangType = "source"
const targetLangType = "target"

type DeepLAdapter struct {
	client *deepl.Client
}

func NewDeepLAdapter(client *deepl.Client) *DeepLAdapter {
	return &DeepLAdapter{client: client}
}

var (
	_ domain.TextTranslator     = (*DeepLAdapter)(nil)
	_ domain.LanguageLister     = (*DeepLAdapter)(nil)
	_ domain.TextRephraser      = (*DeepLAdapter)(nil)
	_ domain.DocumentTranslator = (*DeepLAdapter)(nil)
)

func (a *DeepLAdapter) ListSourceLanguages(ctx context.Context) ([]domain.Language, error) {
	return a.listLanguages(ctx, sourceLangType)
}

func (a *DeepLAdapter) ListTargetLanguages(ctx context.Context) ([]domain.Language, error) {
	return a.listLanguages(ctx, targetLangType)
}

// 列出 DeepL 支持的语言列表 (source 或 target)
func (a *DeepLAdapter) listLanguages(ctx context.Context, langType string) ([]domain.Language, error) {
	query := url.Values{"type": {langType}}
	body, err := a.client.Get(ctx, languagesPath, query)
	if err != nil {
		return nil, err
	}

	var raw []deepl.LanguageDTO
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	langs := make([]domain.Language, len(raw))
	for i, item := range raw {
		langs[i] = domain.Language{
			Code: strings.ToLower(item.Language),
			Name: item.Name,
		}
	}
	return langs, nil
}

// 翻译文本
func (a *DeepLAdapter) Translate(ctx context.Context, entity domain.TranslateEntity) (*domain.TranslateResult, error) {
	if err := entity.Validate(); err != nil {
		return nil, err
	}
	body, err := a.client.PostForm(ctx, translatePath, entity.ToForm())
	if err != nil {
		return nil, err
	}

	var result deepl.TranslateResponseDTO
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if len(result.Translations) == 0 {
		return nil, errors.New("deepl API returned no translations")
	}

	translation := result.Translations[0]
	return &domain.TranslateResult{
		Text:               translation.Text,
		DetectedSourceLang: strings.ToLower(translation.DetectedSourceLanguage),
	}, nil
}

// 改写文本
func (a *DeepLAdapter) Rephrase(ctx context.Context, entity domain.RephraseEntity) (string, error) {
	if err := entity.Validate(); err != nil {
		return "", err
	}

	payload := map[string]any{"text": []string{entity.Text}}
	if entity.Style != "" {
		payload["writing_style"] = entity.Style
	}
	if entity.Tone != "" {
		payload["tone"] = entity.Tone
	}

	body, err := a.client.PostJSON(ctx, rephrasePath, payload)
	if err != nil {
		return "", err
	}

	var result deepl.RephraseResponseDTO
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if len(result.Improvements) == 0 {
		return "", errors.New("deepl API returned no improvements")
	}
	return result.Improvements[0].Text, nil
}

// 翻译文档
func (a *DeepLAdapter) TranslateDocument(ctx context.Context, entity domain.TranslateDocumentEntity) (*deepl.DocumentTranslateDTO, error) {
	// 解析输出文件路径
	outputFile := entity.ResolveOutputFile()
	fields := map[string]string{
		"target_lang": strings.ToUpper(entity.TargetLangCode),
	}
	if entity.SourceLangCode != "" {
		fields["source_lang"] = strings.ToUpper(entity.SourceLangCode)
	}
	if entity.Formality != "" {
		fields["formality"] = entity.Formality
	}
	if entity.GlossaryID != "" {
		fields["glossary_id"] = entity.GlossaryID
	}

	uploadBody, err := a.client.PostMultipart(ctx, documentPath, fields, "file", entity.InputFile)
	if err != nil {
		return nil, err
	}

	var upload deepl.DocumentUploadDTO
	if err := json.Unmarshal(uploadBody, &upload); err != nil {
		return nil, err
	}

	// 获取文档翻译状态
	statusPath := fmt.Sprintf("%s/%s", documentPath, upload.DocumentID)
	form := url.Values{"document_key": {upload.DocumentKey}}

	var status deepl.DocumentStatusDTO
	deadline := time.Now().Add(5 * time.Minute)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if time.Now().After(deadline) {
			return nil, errors.New("document translation timed out")
		}

		statusBody, err := a.client.PostFormPath(ctx, statusPath, form)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(statusBody, &status); err != nil {
			return nil, err
		}

		switch status.Status {
		case "done":
			goto download
		case "error":
			msg := status.ErrorMessage
			if msg == "" {
				msg = status.Message
			}
			if msg == "" {
				msg = "unknown error"
			}
			return nil, fmt.Errorf("document translation failed: %s", msg)
		default:
			time.Sleep(time.Second)
		}
	}

download:
	resultPath := fmt.Sprintf("%s/%s/result", documentPath, upload.DocumentID)
	fileData, err := a.client.PostFormPath(ctx, resultPath, form)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(outputFile, fileData, 0o644); err != nil {
		return nil, err
	}

	// 返回翻译文档结果
	return &deepl.DocumentTranslateDTO{
		Status:           status.Status,
		BilledCharacters: status.BilledCharacters,
		OutputFile:       outputFile,
		TargetLangCode:   strings.ToLower(entity.TargetLangCode),
	}, nil
}
