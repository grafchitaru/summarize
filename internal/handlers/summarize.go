package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/grafchitaru/summarize/internal/domain"
	"github.com/grafchitaru/summarize/internal/middlewares/auth"
	"github.com/grafchitaru/summarize/internal/models"
	"io"
	"net/http"
)

type Sum struct {
	Text string `json:"text"`
}

func (ctx *Handlers) Summarize(res http.ResponseWriter, req *http.Request) {
	var reader io.Reader

	if req.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = req.Body
	}

	body, ioError := io.ReadAll(reader)
	if ioError != nil {
		http.Error(res, ioError.Error(), http.StatusBadRequest)
		return
	}

	var sum Sum

	if err := json.Unmarshal(body, &sum); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	text := sum.Text

	userID, err := auth.GetUserID(req, ctx.Config.SecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	summarizeID := uuid.New()

	summarize := models.NewSummarize{
		Id:     summarizeID.String(),
		UserId: userID,
		Text:   text,
		Status: "Init",
		Tokens: ctx.Ai.GetCountTokens(text),
	}

	err = ctx.Repos.CreateSummarize(summarize)
	if err != nil {
		http.Error(res, "Error Create Summarize:"+fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	result := Result{
		Id: summarizeID.String(),
	}
	data, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(data)

	summarizer := domain.Sum{
		Id:     summarizeID.String(),
		Text:   text,
		Prompt: ctx.Config.AiSummarizePrompt,
		Ai:     ctx.Ai,
		Repos:  ctx.Repos,
	}
	domain.Summarize(summarizer)
}
