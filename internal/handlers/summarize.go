package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/grafchitaru/summarize/internal/middlewares/auth"
	"github.com/grafchitaru/summarize/internal/models"
	"io"
	"net/http"
	"sync"
)

type Sum struct {
	Text string `json:"text"`
}

func (ctx *HandlerContext) Summarize(res http.ResponseWriter, req *http.Request) {
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

	chunks := SplitTextIntoChunks(text, 20000)

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

	go func() {
		var wg sync.WaitGroup
		var summarizedTexts []string

		wg.Add(len(chunks))

		for _, chunk := range chunks {
			go func(chunk string) {
				defer wg.Done()

				summarizedChunk, err := ctx.Ai.Send(chunk, ctx.Config.AiSummarizePrompt)
				if err != nil {
					http.Error(res, err.Error(), http.StatusBadRequest)
					return
				}

				summarizedTexts = append(summarizedTexts, summarizedChunk)
			}(chunk)
		}

		wg.Wait()

		finalSummarizedText := fmt.Sprintf("%s", summarizedTexts)

		finalSummarizedText, err = ctx.Ai.Send(finalSummarizedText, ctx.Config.AiSummarizePrompt)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		err = ctx.Repos.UpdateSummarizeResult(summarizeID.String(), "Complete", finalSummarizedText)
		if err != nil {
			http.Error(res, "Error Save Summarize:"+fmt.Sprintf("%d", err), http.StatusInternalServerError)
			return
		}
	}()
}

func SplitTextIntoChunks(text string, chunkSize int) []string {
	var chunks []string
	runes := []rune(text)
	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}
	return chunks
}
