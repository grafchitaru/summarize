package gigachat

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
	"unicode"
)

const baseUrl = "https://gigachat.devices.sberbank.ru/api/v1"
const authUrl = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
const scope = "scope=GIGACHAT_API_PERS"

type Ai struct {
	AccessToken string
}

type Payload struct {
	Model             string    `json:"model"`
	Messages          []Message `json:"messages"`
	Temperature       float64   `json:"temperature"`
	TopP              float64   `json:"top_p"`
	N                 int       `json:"n"`
	Stream            bool      `json:"stream"`
	MaxTokens         int       `json:"max_tokens"`
	RepetitionPenalty float64   `json:"repetition_penalty"`
	UpdateInterval    int       `json:"update_interval"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type Completions struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
		SystemTokens     int `json:"system_tokens"`
	} `json:"usage"`
}

func New(clientId string, clientSecret string) (*Ai, error) {
	const op = "ai.gigachat.New"

	jsonData, err := getAccessToken(clientId, clientSecret)
	if err != nil {
		return nil, fmt.Errorf("%s: error create access token: %w", op, err)
	}

	var tokenResponse TokenResponse
	err = json.Unmarshal([]byte(jsonData), &tokenResponse)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(tokenResponse.AccessToken)
	return &Ai{AccessToken: tokenResponse.AccessToken}, nil
}

func isString(s interface{}) bool {
	_, ok := s.(string)
	return ok
}

func getAccessToken(clientId string, clientSecret string) (string, error) {
	const op = "ai.gigachat.getAccessToken"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}
	req, err := http.NewRequest("POST", authUrl, strings.NewReader(scope))
	if err != nil {
		log.Printf("%s: error create req: %v", op, err) // Логирование ошибки
		return "", fmt.Errorf("%s: error create req: %w", op, err)
	}

	combined := clientId + ":" + clientSecret
	token := base64.StdEncoding.EncodeToString([]byte(combined))

	RqUID := uuid.New().String()
	req.Header.Add("RqUID", RqUID)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+token)

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: error send req: %w", op, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("%s: error get res: %w", op, err)
	}

	return string(body), nil
}

func (s *Ai) Send(text string, prompt string) (string, error) {
	const op = "ai.gigachat.Send"
	payload := Payload{
		Model:             "GigaChat",
		Messages:          []Message{{Role: "user", Content: prompt + text}},
		Temperature:       1,
		TopP:              0.1,
		N:                 1,
		Stream:            false,
		MaxTokens:         512,
		RepetitionPenalty: 1,
		UpdateInterval:    0,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("%s: error Marshal: %w", op, err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", baseUrl+"/chat/completions", bytes.NewReader(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("%s: error create req: %w", op, err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: error send: %w", op, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: unexpected status code: %d", op, res.StatusCode)
	}

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("%s: error read res: %w", op, err)
	}

	var completions Completions
	err = json.Unmarshal([]byte(jsonData), &completions)
	if err != nil {
		return "", fmt.Errorf("%s: error unmarshal res: %w", op, err)
	}

	if len(completions.Choices) == 0 {
		return "", fmt.Errorf("%s: no completions found", op)
	}

	return completions.Choices[0].Message.Content, nil
}

func (s *Ai) GetCountTokens(text string) int {
	var tokenLength int
	latinTokens := 0
	unicodeTokens := 0

	for _, r := range text {
		if unicode.IsLetter(r) {
			if unicode.Is(unicode.Latin, r) {
				latinTokens++
			} else if unicode.Is(unicode.Cyrillic, r) {
				unicodeTokens++
			}
		}
	}

	if latinTokens > unicodeTokens {
		tokenLength = 4
	} else if unicodeTokens > latinTokens {
		tokenLength = 2
	} else {
		fmt.Println("The language of the text cannot be determined")
		tokenLength = 4
	}

	tokens := int(math.Ceil(float64(len(text)) / float64(tokenLength)))

	return tokens
}
