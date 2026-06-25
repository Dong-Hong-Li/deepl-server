package delivery

import (
	"context"
	"deep-map-server/config"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type HTTPRouterRegistry interface {
	RegisterRouter(r chi.Router)
}

func RunHTTPServer(ctx context.Context, cfg *config.Config, server *mcp.Server, routerRegistries []HTTPRouterRegistry) error {
	handler := NewHTTPHandler(cfg.AuthToken(), server, routerRegistries)

	httpServer := &http.Server{
		Addr:    cfg.HTTPAddr(),
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutdownCtx)
	}()

	log.Printf("DeepL HTTP server listening on http://%s (MCP: /mcp, API: /api/*)", cfg.HTTPAddr())
	if cfg.AuthToken() != "" {
		log.Print("MCP auth enabled (Authorization: Bearer MCP_AUTH_TOKEN)")
	}

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server: %w", err)
	}
	return nil
}

func NewHTTPHandler(authToken string, server *mcp.Server, routerRegistries []HTTPRouterRegistry) http.Handler {
	mcpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)

	r := chi.NewRouter()
	r.Use(corsMiddleware)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Group(func(r chi.Router) {
		if authToken != "" {
			r.Use(bearerAuth(authToken))
		}
		r.Handle("/mcp", mcpHandler)
		for _, registry := range routerRegistries {
			registry.RegisterRouter(r)
		}
	})

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func bearerAuth(expected string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			const prefix = "Bearer "
			if !strings.HasPrefix(auth, prefix) || strings.TrimPrefix(auth, prefix) != expected {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
