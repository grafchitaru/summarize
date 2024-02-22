package main

import (
	"fmt"
	"github.com/grafchitaru/summarize/internal/config"
	"github.com/grafchitaru/summarize/internal/server"
	storage2 "github.com/grafchitaru/summarize/internal/storage"
	"github.com/grafchitaru/summarize/internal/storage/postgresql"
)

func main() {
	cfg := *config.NewConfig()

	var storage storage2.Repositories
	var err error

	storage, err = postgresql.New(cfg.PostgresDatabaseDsn)
	if err != nil {
		fmt.Println("Error initialize storage: %w", err)
	}

	defer storage.Close()

	server.New(config.HandlerContext{Config: cfg, Repos: storage})
}
