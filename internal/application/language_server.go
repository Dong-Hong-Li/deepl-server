package application

import (
	"context"
	"deep-map-server/internal/domain"
)

type LanguageServer struct {
	lister domain.LanguageLister
}

func NewLanguageServer(lister domain.LanguageLister) *LanguageServer {
	return &LanguageServer{lister: lister}
}

// 获取可供翻译的源语言列表
func (s *LanguageServer) GetSourceLanguages(ctx context.Context) ([]domain.Language, error) {
	return s.lister.ListSourceLanguages(ctx)
}

// 获取可供翻译的目标语言列表
func (s *LanguageServer) GetTargetLanguages(ctx context.Context) ([]domain.Language, error) {
	return s.lister.ListTargetLanguages(ctx)
}
