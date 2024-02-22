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

	err := http.ListenAndServe(ctx.Config.HTTPServerAddress, r)
	if err != nil {
		fmt.Println("Error server: %w", err)
	}
}
