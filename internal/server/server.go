package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/summarize/internal/handlers"
	"github.com/grafchitaru/summarize/internal/middlewares/auth"
	"github.com/grafchitaru/summarize/internal/middlewares/compress"
	"github.com/grafchitaru/summarize/internal/middlewares/logger"
	"net/http"
)

func New(ctx handlers.HandlerContext) {
	hc := &handlers.HandlerContext{
		Config: ctx.Config,
		Repos:  ctx.Repos,
		Ai:     ctx.Ai,
	}

	r := chi.NewRouter()

	r.Use(logger.WithLogging)
	r.Use(compress.WithCompressionResponse)
	r.Use(auth.WithUserCookie(hc.Config.SecretKey))

	r.Post("/ping", hc.Ping)

	r.Post("/api/user/register", hc.Register)

	r.Post("/api/user/login", hc.Login)

	r.Post("/api/user/summarize", hc.Summarize)

	r.Get("/api/user/summarize/{id}", hc.GetSummarizeText)

	r.Get("/api/user/stat", hc.Stat)

	r.Get("/api/user/status", hc.Status)

	err := http.ListenAndServe(ctx.Config.HTTPServerAddress, r)
	if err != nil {
		fmt.Println("Error server: %w", err)
	}
}
