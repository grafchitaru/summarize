package handlers

import (
	ai2 "github.com/grafchitaru/summarize/internal/ai"
	"github.com/grafchitaru/summarize/internal/auth"
	"github.com/grafchitaru/summarize/internal/config"
	"github.com/grafchitaru/summarize/internal/storage"
)

type Handlers struct {
	Config config.Config
	Repos  storage.Repositories
	Ai     ai2.AI
	Auth   auth.AuthService
}
