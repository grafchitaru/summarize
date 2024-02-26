package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/summarize/internal/config"
	"github.com/grafchitaru/summarize/internal/handlers"
	"github.com/grafchitaru/summarize/internal/middlewares/auth"
	"github.com/grafchitaru/summarize/internal/middlewares/compress"
	"github.com/grafchitaru/summarize/internal/middlewares/logger"
	"net/http"
)

func New(ctx config.HandlerContext) {

	r := chi.NewRouter()

	r.Use(logger.WithLogging)
	r.Use(compress.WithCompressionResponse)
	r.Use(auth.WithUserCookie(ctx))

	r.Get("/ping", func(res http.ResponseWriter, req *http.Request) {
		handlers.Ping(ctx, res)
	})

	r.Post("/api/user/register", func(res http.ResponseWriter, req *http.Request) {
		handlers.Register(ctx, res, req)
	})

	r.Post("/api/user/login", func(res http.ResponseWriter, req *http.Request) {
		handlers.Login(ctx, res, req)
	})

	r.Post("/api/user/summarize", func(res http.ResponseWriter, req *http.Request) {
		handlers.Summarize(ctx, res, req)
	})

	r.Get("/api/user/summarize/{id}", func(res http.ResponseWriter, req *http.Request) {
		handlers.GetSummarizeText(ctx, res, req)
	})

	r.Get("/api/user/stat", func(res http.ResponseWriter, req *http.Request) {
		handlers.Stat(ctx, res)
	})

	r.Get("/api/user/status", func(res http.ResponseWriter, req *http.Request) {
		handlers.Status(ctx, res)
	})

	err := http.ListenAndServe(ctx.Config.HTTPServerAddress, r)
	if err != nil {
		fmt.Println("Error server: %w", err)
	}
}
