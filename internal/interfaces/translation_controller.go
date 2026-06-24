package interfaces

import (
	"deep-map-server/internal/application"
	"deep-map-server/internal/delivery/binding"

	"github.com/go-chi/chi/v5"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TranslationController struct {
	TranslateTextServer     *application.TranslateTextServer
	LanguageServer          *application.LanguageServer
	RephraseTextServer      *application.RephraseTextServer
	DocumentTranslateServer *application.DocumentTranslateServer
}

func NewTranslationController(
	translateTextServer *application.TranslateTextServer,
	languageServer *application.LanguageServer,
	rephraseTextServer *application.RephraseTextServer,
	documentTranslateServer *application.DocumentTranslateServer,
) *TranslationController {
	return &TranslationController{
		TranslateTextServer:     translateTextServer,
		LanguageServer:          languageServer,
		RephraseTextServer:      rephraseTextServer,
		DocumentTranslateServer: documentTranslateServer,
	}
}

func (c *TranslationController) RegisterRouter(r chi.Router) {
	r.Get("/api/languages/source", c.handleGetSourceLanguages)
	r.Get("/api/languages/target", c.handleGetTargetLanguages)
	r.Post("/api/translate", c.handleTranslateText)
	r.Get("/api/writing-styles", c.handleGetWritingStyles)
	r.Get("/api/writing-tones", c.handleGetWritingTones)
	r.Post("/api/rephrase", c.handleRephraseText)
	r.Post("/api/translate-document", c.handleTranslateDocument)
}

func (c *TranslationController) RegisterTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-source-languages",
		Description: "Get list of available source languages for translation",
	}, binding.Exec(c.getSourceLanguages))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-target-languages",
		Description: "Get list of available target languages for translation",
	}, binding.Exec(c.getTargetLanguages))

	mcp.AddTool(server, &mcp.Tool{
		Name: "translate-text",
		Description: "Translate text to a target language using DeepL API. " +
			"When using a glossary, sourceLangCode is required.",
	}, binding.Handle(c.translateText))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-writing-styles",
		Description: "Get list of writing styles the DeepL API can use while rephrasing text",
	}, binding.Exec(c.getWritingStyles))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-writing-tones",
		Description: "Get list of writing tones the DeepL API can use while rephrasing text",
	}, binding.Exec(c.getWritingTones))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "rephrase-text",
		Description: "Rephrase text in the same language using DeepL API",
	}, binding.Handle(c.rephraseText))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "translate-document",
		Description: "Translate a document file using DeepL API",
	}, binding.Handle(c.translateDocument))
}
