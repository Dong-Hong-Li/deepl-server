package request

// TranslateDocumentRequest translate-document 工具入参。
type TranslateDocumentRequest struct {
	InputFile      string `json:"inputFile" jsonschema:"Path to the input document file to translate" validate:"required"`
	OutputFile     string `json:"outputFile,omitempty" jsonschema:"Path where the translated document will be saved (auto-generated if omitted)"`
	SourceLangCode string `json:"sourceLangCode,omitempty" jsonschema:"Source language code, or leave empty for auto-detection"`
	TargetLangCode string `json:"targetLangCode" jsonschema:"Target language code" validate:"required"`
	Formality      string `json:"formality,omitempty" jsonschema:"Controls formality: less, more, default, prefer_less, prefer_more"`
	GlossaryID     string `json:"glossaryId,omitempty" jsonschema:"Glossary ID to ensure consistent terminology translation"`
}
