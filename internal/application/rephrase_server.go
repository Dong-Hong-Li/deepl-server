package application

import (
	"context"
	"deep-map-server/internal/domain"
)

type RephraseTextServer struct {
	rephraser domain.TextRephraser
}

func NewRephraseTextServer(rephraser domain.TextRephraser) *RephraseTextServer {
	return &RephraseTextServer{rephraser: rephraser}
}

// 改写文本
func (s *RephraseTextServer) RephraseText(ctx context.Context, entity domain.RephraseEntity) (string, error) {
	return s.rephraser.Rephrase(ctx, entity)
}
