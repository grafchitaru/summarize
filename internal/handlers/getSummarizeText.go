package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/grafchitaru/summarize/internal/config"
	"github.com/grafchitaru/summarize/internal/middlewares/auth"
	"net/http"
)

func GetSummarizeText(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
	summarizeID := chi.URLParam(req, "id")
	if summarizeID == "" {
		http.Error(res, "ID not found", http.StatusNotFound)
		return
	}

	userID, err := auth.GetUserID(req, ctx.Config.SecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	result, err := ctx.Repos.GetSummarize(summarizeID, userID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(data)
}
