package main

import (
	"fmt"
	ai2 "github.com/grafchitaru/summarize/internal/ai"
	"github.com/grafchitaru/summarize/internal/ai/gigachat"
	"github.com/grafchitaru/summarize/internal/config"
	"github.com/grafchitaru/summarize/internal/handlers"
	"github.com/grafchitaru/summarize/internal/server"
	storage2 "github.com/grafchitaru/summarize/internal/storage"
	"github.com/grafchitaru/summarize/internal/storage/postgresql"
)

func main() {
	cfg := *config.NewConfig()

	var storage storage2.Repositories
	var err error
	var ai ai2.AI

	storage, err = postgresql.New(cfg.PostgresDatabaseDsn)
	if err != nil {
		fmt.Println("Error initialize storage: %w", err)
	}

	ai, err = gigachat.New(cfg.GigaChatClientId, cfg.GigaChatClientSecret)
	if err != nil {
		fmt.Println("Error initialize ai: %w", err)
	}

	defer storage.Close()

	server.New(handlers.HandlerContext{Config: cfg, Repos: storage, Ai: ai})
}
