package application

import (
	"context"
	"crypto/sha256"
	"deep-map-server/internal/delivery"
	"deep-map-server/internal/domain"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
)

type DocumentTranslateServer struct {
	translator    domain.DocumentTranslator
	maxFileSize   int
	objectStorage domain.ObjectStorage
}

func NewDocumentTranslateServer(translator domain.DocumentTranslator, maxFileSize int, objectStorage domain.ObjectStorage) *DocumentTranslateServer {
	return &DocumentTranslateServer{
		translator:    translator,
		maxFileSize:   maxFileSize,
		objectStorage: objectStorage,
	}
}

// 翻译文档
func (s *DocumentTranslateServer) TranslateDocument(ctx context.Context, entity domain.TranslateDocumentEntity) (*domain.TranslateDocumentResult, error) {
	if err := entity.Validate(); err != nil {
		return nil, err
	}
	if err := domain.ValidateDocumentInputFile(entity.InputFile, s.maxFileSize); err != nil {
		return nil, err
	}
	content, err := os.ReadFile(entity.InputFile)
	if err != nil {
		return nil, err
	}
	// TODO: 已存在则直接复用，不再重复翻译与落库
	hash := hashResumeContent(content)
	if hash == "" {
		return nil, errors.New("hash is empty")
	}
	entity.OutputFile = hash + filepath.Ext(entity.InputFile)

	// TODO: 先落库，再翻译
	// InsertResume

	// 将源文件添加进对象存储
	fileKey := delivery.NewFileKey("before_translate", filepath.Ext(entity.InputFile))
	s.objectStorage.Upload(ctx, fileKey, content, "application/octet-stream")

	result, err := s.translator.TranslateDocument(ctx, entity)
	if err != nil {
		s.objectStorage.Delete(ctx, entity.InputFile)
		return nil, err
	}
	// 将翻译后的文件添加进对象存储
	content, err = os.ReadFile(result.OutputFile)
	if err != nil {
		return nil, err
	}

	fileKey = delivery.NewFileKey("after_translate", filepath.Ext(entity.InputFile))
	s.objectStorage.Upload(ctx, fileKey, content, "application/octet-stream")

	filePath, err := s.objectStorage.PresignGetObjectURL(ctx, fileKey)
	if err != nil {
		return nil, err
	}
	return &domain.TranslateDocumentResult{
		Status:           result.Status,
		BilledCharacters: result.BilledCharacters,
		OutputFilePath:   filePath,
		TargetLangCode:   result.TargetLangCode,
	}, nil
}

// 计算简历内容的 SHA-256 哈希值
func hashResumeContent(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
