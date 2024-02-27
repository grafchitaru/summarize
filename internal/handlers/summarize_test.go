package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/grafchitaru/summarize/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafchitaru/summarize/internal/mocks"
)

type res struct {
	Id string `json:"id"`
}

func TestSummarize(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "556501c0-a97d-47ed-9add-73b4a4116c83"
	mockStorage := &mocks.MockStorage{
		CreateSummarizeFunc: func(id string, user_id string, text string, status string, tokens int) error {
			return nil
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	ctx := config.HandlerContext{Config: *cfg, Repos: mockStorage, Auth: mockAuthService}

	body, _ := json.Marshal(Sum{Text: "В значительной степени обуславливает создание направлений прогрессивного развития."})
	req, err := http.NewRequest("POST", "/api/user/summarize", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)
	r := httptest.NewRecorder()

	Summarize(ctx, r, req)

	rr := httptest.NewRecorder()

	assert.Equal(t, rr.Code, http.StatusOK)

	var res res
	json.NewDecoder(r.Body).Decode(&res)
}

func TestSummarize_CreateError(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "556501c0-a97d-47ed-9add-73b4a4116c83"
	mockStorage := &mocks.MockStorage{
		CreateSummarizeFunc: func(id string, user_id string, text string, status string, tokens int) error {
			return errors.New("create summarize error")
		},
	}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	ctx := config.HandlerContext{Config: *cfg, Repos: mockStorage, Auth: mockAuthService}

	body, _ := json.Marshal(Sum{Text: "В значительной степени обуславливает создание направлений прогрессивного развития."})
	req, err := http.NewRequest("POST", "/api/user/summarize", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)
	r := httptest.NewRecorder()

	Summarize(ctx, r, req)

	assert.Equal(t, http.StatusUnauthorized, r.Code)
}
