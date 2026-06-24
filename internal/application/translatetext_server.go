package application

import (
	"context"
	"deep-map-server/internal/domain"
)

type TranslateTextServer struct {
	translator domain.TextTranslator
}

func NewTranslateTextServer(translator domain.TextTranslator) *TranslateTextServer {
	return &TranslateTextServer{translator: translator}
}

// 翻译文本
func (s *TranslateTextServer) TranslateText(ctx context.Context, entity domain.TranslateEntity) (*domain.TranslateResult, error) {
	// TODO: 数据库查询翻译历史，如果有则直接返回，如果没有则翻译并保存到数据库
	return s.translator.Translate(ctx, entity)
}
