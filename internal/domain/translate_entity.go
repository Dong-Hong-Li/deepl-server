package domain

import (
	"errors"
	"net/url"
	"strings"
)

type TranslateEntity struct {
	SourceLangCode string
	TargetLangCode string
	Formality      string
	GlossaryID     string
	Text           string
}

func (e *TranslateEntity) Validate() error {
	if strings.TrimSpace(e.Text) == "" {
		return errors.New("text is required")
	}
	if strings.TrimSpace(e.TargetLangCode) == "" {
		return errors.New("targetLangCode is required")
	}
	return nil
}

func (e *TranslateEntity) ToForm() url.Values {
	form := url.Values{}
	form.Set("text", e.Text)
	form.Set("target_lang", strings.ToUpper(e.TargetLangCode))
	if e.SourceLangCode != "" {
		form.Set("source_lang", strings.ToUpper(e.SourceLangCode))
	}
	if e.Formality != "" {
		form.Set("formality", e.Formality)
	}
	if e.GlossaryID != "" {
		form.Set("glossary_id", e.GlossaryID)
	}
	return form
}
