package deepl

type LanguageDTO struct {
	Language string `json:"language"`
	Name     string `json:"name"`
}

type TranslateResponseDTO struct {
	Translations []struct {
		Text                   string `json:"text"`
		DetectedSourceLanguage string `json:"detected_source_language"`
	} `json:"translations"`
}

type RephraseResponseDTO struct {
	Improvements []struct {
		Text string `json:"text"`
	} `json:"improvements"`
}

type DocumentUploadDTO struct {
	DocumentID  string `json:"document_id"`
	DocumentKey string `json:"document_key"`
}

type DocumentStatusDTO struct {
	Status           string `json:"status"`
	BilledCharacters int    `json:"billed_characters"`
	ErrorMessage     string `json:"error_message"`
	Message          string `json:"message"`
}

type DocumentTranslateDTO struct {
	Status           string `json:"status"`
	BilledCharacters int    `json:"billed_characters"`
	OutputFile       string `json:"output_file"`
	TargetLangCode   string `json:"target_lang_code"`
}
