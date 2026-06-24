package domain

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type TranslateDocumentEntity struct {
	// 输入文件
	InputFile string
	// 输出文件
	OutputFile string
	// 源语言代码
	SourceLangCode string
	// 目标语言代码
	TargetLangCode string
	// 语体
	Formality string
	// 语料库ID
	GlossaryID string
}

type TranslateDocumentResult struct {
	// 翻译状态
	Status string
	// 计费字符数
	BilledCharacters int
	// 输出文件路径
	OutputFilePath string
	// 目标语言代码
	TargetLangCode string
}

func (e *TranslateDocumentEntity) Validate() error {
	if strings.TrimSpace(e.InputFile) == "" {
		return errors.New("inputFile is required")
	}
	if strings.TrimSpace(e.TargetLangCode) == "" {
		return errors.New("targetLangCode is required")
	}
	if err := ValidateDocumentInputFileType(e.InputFile); err != nil {
		return err
	}
	return nil
}

func ValidateDocumentInputFile(inputFile string, maxBytes int) error {
	info, err := os.Stat(inputFile)
	if err != nil {
		return err
	}
	if maxBytes > 0 && info.Size() > int64(maxBytes) {
		return fmt.Errorf("inputFile exceeds max size limit (%d bytes)", maxBytes)
	}
	return nil
}

// SupportedDocumentExtensions DeepL 文档翻译 API 支持的扩展名。
// 见 https://developers.deepl.com/docs/api-reference/document/upload-and-translate-a-document
var SupportedDocumentExtensions = map[string]struct{}{
	".docx":  {},
	".pptx":  {},
	".xlsx":  {},
	".pdf":   {},
	".htm":   {},
	".html":  {},
	".txt":   {},
	".xlf":   {},
	".xliff": {},
	".srt":   {},
	".jpeg":  {},
	".jpg":   {},
	".png":   {},
}

func ValidateDocumentInputFileType(inputFile string) error {
	ext := strings.ToLower(filepath.Ext(inputFile))
	if _, ok := SupportedDocumentExtensions[ext]; !ok {
		return errors.New("inputFile type is not supported by DeepL document translation API; supported: docx, pptx, xlsx, pdf, htm, html, txt, xlf, xliff, srt, jpeg, jpg, png")
	}
	return nil
}

// ResolveOutputFile 若未指定输出路径，按官方规则自动生成（如 doc_de.pdf）。
func (e *TranslateDocumentEntity) ResolveOutputFile() string {
	if e.OutputFile != "" {
		return e.OutputFile
	}

	dir := filepath.Dir(e.InputFile)
	base := filepath.Base(e.InputFile)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	langCode := strings.ToLower(e.TargetLangCode)
	if idx := strings.Index(langCode, "-"); idx > 0 {
		langCode = langCode[:idx]
	}

	return filepath.Join(dir, name+"_"+langCode+ext)
}
