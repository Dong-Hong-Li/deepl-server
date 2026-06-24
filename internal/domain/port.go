package domain

import (
	"context"
	"deep-map-server/internal/infrastructure/deepl"
)

// TextTranslator 文本翻译端口。
type TextTranslator interface {
	Translate(ctx context.Context, entity TranslateEntity) (*TranslateResult, error)
}

// LanguageLister 语言列表端口。
type LanguageLister interface {
	ListSourceLanguages(ctx context.Context) ([]Language, error)
	ListTargetLanguages(ctx context.Context) ([]Language, error)
}

// TextRephraser 文本改写端口。
type TextRephraser interface {
	Rephrase(ctx context.Context, entity RephraseEntity) (string, error)
}

// DocumentTranslator 文档翻译端口。
type DocumentTranslator interface {
	TranslateDocument(ctx context.Context, entity TranslateDocumentEntity) (*deepl.DocumentTranslateDTO, error)
}

type ObjectStorage interface {
	// 上传文件
	Upload(ctx context.Context, fileKey string, file []byte, contentType string) error
	// 下载文件
	Download(ctx context.Context, fileKey string) ([]byte, error)
	// 删除文件
	Delete(ctx context.Context, fileKey string) error
	// 获取文件URL
	PresignGetObjectURL(ctx context.Context, objectKey string) (string, error)
}
