package domain

import (
	"github.com/grafchitaru/summarize/internal/ai"
	"github.com/grafchitaru/summarize/internal/storage"
	"log"
	"strings"
	"sync"
)

type Sum struct {
	Id     string
	Text   string
	Prompt string
	Ai     ai.AI
	Repos  storage.Repositories
}

func Summarize(sum Sum) {
	chunks := SplitTextIntoChunks(sum.Text, 20000)
	var sb strings.Builder

	go func() {
		var wg sync.WaitGroup
		wg.Add(len(chunks))
		//TODO Worker Pool
		for _, chunk := range chunks {
			go func(chunk string) {
				defer wg.Done()

				summarizedChunk, err := sum.Ai.Send(chunk, sum.Prompt)
				if err != nil {
					log.Printf("%s: error summarizedChunk: %v", "handler.summarize", err)
					return
				}
				sb.WriteString(summarizedChunk)
			}(chunk)
		}

		wg.Wait()

		var err error

		finalSummarizedText := sb.String()

		finalSummarizedText, err = sum.Ai.Send(finalSummarizedText, sum.Prompt)
		if err != nil {
			log.Printf("%s: error finalSummarizedText: %v", "handler.summarize", err)
			return
		}

		err = sum.Repos.UpdateSummarizeResult(sum.Id, "Complete", finalSummarizedText)
		if err != nil {
			log.Printf("%s: error UpdateSummarizeResult: %v", "handler.summarize", err)
			return
		}
	}()
	//TODO добавить механизм ретраев
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
