//go:generate go run -mod=mod github.com/google/wire/cmd/wire

//go:build wireinject
// +build wireinject

// wireinject
package main

import (
	"context"
	"deep-map-server/config"
	"deep-map-server/internal/adapter"
	"deep-map-server/internal/application"
	"deep-map-server/internal/domain"
	"deep-map-server/internal/infrastructure/deepl"
	objectstorage "deep-map-server/internal/infrastructure/object_storage"
	"deep-map-server/internal/interfaces"

	"github.com/google/wire"
)

var infrastructureModuleSet = wire.NewSet(
	provideDeepLClient,
	provideObjectStorageClient,
)

var translationModuleSet = wire.NewSet(
	adapter.NewDeepLAdapter,
	adapter.NewObjectOSSAdapter,
	wire.Bind(new(domain.TextTranslator), new(*adapter.DeepLAdapter)),
	wire.Bind(new(domain.LanguageLister), new(*adapter.DeepLAdapter)),
	wire.Bind(new(domain.TextRephraser), new(*adapter.DeepLAdapter)),
	wire.Bind(new(domain.DocumentTranslator), new(*adapter.DeepLAdapter)),
	wire.Bind(new(domain.ObjectStorage), new(*adapter.ObjectOSSAdapter)),
	application.NewTranslateTextServer,
	application.NewLanguageServer,
	application.NewRephraseTextServer,
	provideDocumentTranslateServer,
	provideTranslationController,
)

func provideDeepLClient(cfg *config.Config) *deepl.Client {
	return deepl.NewClient(cfg.GetDeelApiKey())
}

func provideObjectStorageClient(cfg *config.Config) (*objectstorage.Client, error) {
	return objectstorage.NewClient(cfg.OSSConfig(), context.Background())
}

func provideDocumentTranslateServer(
	translator domain.DocumentTranslator,
	cfg *config.Config,
	storage domain.ObjectStorage,
) *application.DocumentTranslateServer {
	return application.NewDocumentTranslateServer(translator, cfg.OSSConfig().MaxFileSize, storage)
}

func provideTranslationController(
	translateTextServer *application.TranslateTextServer,
	languageServer *application.LanguageServer,
	rephraseTextServer *application.RephraseTextServer,
	documentTranslateServer *application.DocumentTranslateServer,
) *interfaces.TranslationController {
	return interfaces.NewTranslationController(
		translateTextServer,
		languageServer,
		rephraseTextServer,
		documentTranslateServer,
	)
}

func initializeApp(cfg *config.Config) (*interfaces.TranslationController, error) {
	panic(wire.Build(
		infrastructureModuleSet,
		translationModuleSet,
	))
}
