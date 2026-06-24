package request

// RephraseTextRequest rephrase-text 工具入参。
type RephraseTextRequest struct {
	Text  string `json:"text" jsonschema:"Text to rephrase" validate:"required"`
	Style string `json:"style,omitempty" jsonschema:"Writing style: academic, business, casual, simple, default, or prefer_* variants"`
	Tone  string `json:"tone,omitempty" jsonschema:"Writing tone: confident, diplomatic, enthusiastic, friendly, default, or prefer_* variants"`
}
