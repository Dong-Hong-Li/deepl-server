package main

import (
	"context"
	"deep-map-server/config"
	"deep-map-server/internal/delivery"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.NewConfig()

	cleanups := make([]func(), 0, 4)
	// TODO: db, err := openDatabase(cfg)
	// cleanups = append(cleanups, func() { _ = db.Close() })
	// TODO: cache, err := openCache(cfg)
	// cleanups = append(cleanups, func() { _ = cache.Close() })

	translation, err := initializeApp(cfg)
	if err != nil {
		log.Fatalf("initialize app: %v", err)
	}
	defer runCleanups(cleanups...)

	server := delivery.NewServertartMCPServer(cfg, ctx, []delivery.MCPToolRegistry{translation})

	switch cfg.Transport() {
	case config.TransportHTTP:
		if err := delivery.RunHTTPServer(ctx, cfg, server, []delivery.HTTPRouterRegistry{translation}); err != nil && err != context.Canceled {
			log.Fatal(err)
		}
	default:
		if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil && err != context.Canceled {
			log.Fatal(err)
		}
	}
}

// runCleanups 按注册相反顺序执行回收（后注册的先关闭，类似 defer 栈）。
func runCleanups(cleanups ...func()) {
	for i := len(cleanups) - 1; i >= 0; i-- {
		if cleanups[i] != nil {
			cleanups[i]()
		}
	}
}
