package domain

import (
	"errors"
	"strings"
)

type RephraseEntity struct {
	Text  string
	Style string
	Tone  string
}

func (e *RephraseEntity) Validate() error {
	if strings.TrimSpace(e.Text) == "" {
		return errors.New("text is required")
	}
	if e.Style != "" && e.Tone != "" {
		return errors.New("style and tone cannot both be set")
	}
	return nil
}
