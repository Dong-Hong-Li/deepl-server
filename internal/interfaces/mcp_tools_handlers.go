package interfaces

import (
	"context"
	"deep-map-server/internal/domain"
	"deep-map-server/internal/interfaces/request"
	"deep-map-server/internal/interfaces/response"
	"fmt"
)

// -------------------------------- MCP Tools --------------------------------

/**
* @description: 获取可供翻译的源语言列表
* @param {context.Context} ctx
* @return {*response.LanguageListResponse, error}
 */
func (c *TranslationController) getSourceLanguages(ctx context.Context) (*response.LanguageListResponse, error) {
	langs, err := c.LanguageServer.GetSourceLanguages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get source languages: %w", err)
	}
	return response.NewLanguageList(langs), nil
}

/**
* @description: 获取可供翻译的目标语言列表
* @param {context.Context} ctx
* @return {*response.LanguageListResponse, error}
 */
func (c *TranslationController) getTargetLanguages(ctx context.Context) (*response.LanguageListResponse, error) {
	langs, err := c.LanguageServer.GetTargetLanguages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get target languages: %w", err)
	}
	return response.NewLanguageList(langs), nil
}

/**
* @description: 翻译文本
* @param {context.Context} ctx
* @param {request.TranslateTextRequest} req
* @return {*response.TranslateTextResponse, error}
 */
func (c *TranslationController) translateText(ctx context.Context, req request.TranslateTextRequest) (*response.TranslateTextResponse, error) {
	if req.GlossaryID != "" && req.SourceLangCode == "" {
		return nil, fmt.Errorf("sourceLangCode is required when glossaryId is set")
	}

	result, err := c.TranslateTextServer.TranslateText(ctx, domain.TranslateEntity{
		SourceLangCode: req.SourceLangCode,
		TargetLangCode: req.TargetLangCode,
		Formality:      req.Formality,
		GlossaryID:     req.GlossaryID,
		Text:           req.Text,
	})
	if err != nil {
		return nil, fmt.Errorf("translation failed: %w", err)
	}

	return &response.TranslateTextResponse{
		Text:               result.Text,
		DetectedSourceLang: result.DetectedSourceLang,
		TargetLangCode:     req.TargetLangCode,
	}, nil
}

/**
* @description: 获取可供改写的写作风格列表
* @param {context.Context} ctx
* @return {*response.WritingStylesResponse, error}
 */
func (c *TranslationController) getWritingStyles(ctx context.Context) (*response.WritingStylesResponse, error) {
	return response.NewWritingStylesResponse(), nil
}

/**
* @description: 获取可供改写的写作语气列表
* @param {context.Context} ctx
* @return {*response.WritingTonesResponse, error}
 */
func (c *TranslationController) getWritingTones(ctx context.Context) (*response.WritingTonesResponse, error) {
	return response.NewWritingTonesResponse(), nil
}

/**
* @description: 改写文本
* @param {context.Context} ctx
* @param {request.RephraseTextRequest} req
* @return {*response.RephraseTextResponse, error}
 */
func (c *TranslationController) rephraseText(ctx context.Context, req request.RephraseTextRequest) (*response.RephraseTextResponse, error) {
	text, err := c.RephraseTextServer.RephraseText(ctx, domain.RephraseEntity{
		Text:  req.Text,
		Style: req.Style,
		Tone:  req.Tone,
	})
	if err != nil {
		return nil, fmt.Errorf("rephrasing failed: %w", err)
	}
	return &response.RephraseTextResponse{Text: text}, nil
}

/**
* @description: 翻译文档
* @param {context.Context} ctx
* @param {request.TranslateDocumentRequest} req
* @return {*response.TranslateDocumentResponse, error}
 */
func (c *TranslationController) translateDocument(ctx context.Context, req request.TranslateDocumentRequest) (*response.TranslateDocumentResponse, error) {
	result, err := c.DocumentTranslateServer.TranslateDocument(ctx, domain.TranslateDocumentEntity{
		InputFile:      req.InputFile,
		OutputFile:     req.OutputFile,
		SourceLangCode: req.SourceLangCode,
		TargetLangCode: req.TargetLangCode,
		Formality:      req.Formality,
		GlossaryID:     req.GlossaryID,
	})
	if err != nil {
		return nil, fmt.Errorf("document translation failed: %w", err)
	}

	return &response.TranslateDocumentResponse{
		Status:           result.Status,
		BilledCharacters: result.BilledCharacters,
		OutputFilePath:   result.OutputFilePath,
		TargetLangCode:   result.TargetLangCode,
	}, nil
}
