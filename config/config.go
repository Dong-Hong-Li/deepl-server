package config

import (
	"log"
	"os"
)

const (
	TransportStdio = "stdio"
	TransportHTTP  = "http"
)

type Config struct {
	deeplAPIKey string
	transport   string
	httpAddr    string
	authToken   string
	ossConfig   *OSSConfig
}

func NewConfig() *Config {
	deeplAPIKey := os.Getenv("DEEPL_API_KEY")
	if deeplAPIKey == "" {
		log.Fatalf("DEEPL_API_KEY is not set")
	}

	transport := os.Getenv("MCP_TRANSPORT")
	if transport == "" {
		transport = TransportStdio
	}
	if transport != TransportStdio && transport != TransportHTTP {
		log.Fatalf("MCP_TRANSPORT must be %q or %q, got %q", TransportStdio, TransportHTTP, transport)
	}

	httpAddr := os.Getenv("MCP_HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8080"
	}

	ossConfig, err := validateOSSConfig()
	if err != nil {
		log.Fatalf("validate oss config: %v", err)
	}

	return &Config{
		deeplAPIKey: deeplAPIKey,
		transport:   transport,
		httpAddr:    httpAddr,
		authToken:   os.Getenv("MCP_AUTH_TOKEN"),
		ossConfig:   ossConfig,
	}
}

func (c *Config) GetDeelApiKey() string {
	return c.deeplAPIKey
}

func (c *Config) Transport() string {
	return c.transport
}

func (c *Config) HTTPAddr() string {
	return c.httpAddr
}

func (c *Config) AuthToken() string {
	return c.authToken
}

func (c *Config) OSSConfig() *OSSConfig {
	return c.ossConfig
}
