package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	ai2 "github.com/grafchitaru/summarize/internal/ai"
	"github.com/grafchitaru/summarize/internal/auth"
	"github.com/grafchitaru/summarize/internal/storage"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	HTTPServerAddress    string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8119"`
	PostgresDatabaseDsn  string `env:"DATABASE_DSN" envDefault:"postgres://root:root@localhost:54322/app"`
	SecretKey            string `env:"SECRET_KEY" envDefault:"your_secret_key"`
	GigaChatClientId     string `env:"GIGACHAT_CLIENT_ID"`
	GigaChatClientSecret string `env:"GIGACHAT_CLIENT_SECRET"`
	AiMaxLenText         int    `env:"AI_MAX_LEN_TEXT" envDefault:"100000"`
	AiSummarizePrompt    string `env:"AI_SUMMARIZE_PROMPT" envDefault:"Напиши краткое изложение следующего содержания:\n"`
}

type HandlerContext struct {
	Config Config
	Repos  storage.Repositories
	Ai     ai2.AI
	Auth   auth.AuthService
}

type Configs interface {
	NewConfig() *Config
}

func NewConfig() *Config {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = env.Parse(&cfg)
	if err != nil {
		fmt.Println("Can't parse  config: %w", err)
	}

	flag.StringVar(&cfg.HTTPServerAddress, "a", cfg.HTTPServerAddress, "HTTP server address")
	flag.StringVar(&cfg.PostgresDatabaseDsn, "d", cfg.PostgresDatabaseDsn, "PostgreSql database dsn")

	flag.Parse()

	return &cfg
}
