package handlers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafchitaru/summarize/internal/mocks"
)

func TestGetSummarizeText(t *testing.T) {
	cfg := mocks.NewConfig()
	testUserID := "556501c0-a97d-47ed-9add-73b4a4116c83"
	testSummarizeID := "556501c0-a97d-47ed-9add-73b4a4116c83"
	mockStorage := &mocks.MockStorage{}

	mockAuthService := mocks.NewMockAuthService()
	mockAuthService.GetUserIDFunc = func(req *http.Request, secretKey string) (string, error) {
		return testUserID, nil
	}

	req, err := http.NewRequest("GET", "/api/user/summarize/"+testSummarizeID, bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)
	r := httptest.NewRecorder()

	hc := &Handlers{
		Config: *cfg,
		Repos:  mockStorage,
		Auth:   mockAuthService,
	}
	hc.GetSummarizeText(r, req)

	rr := httptest.NewRecorder()

	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestGetSummarizeTextError(t *testing.T) {
	cfg := mocks.NewConfig()
	mockStorage := &mocks.MockStorage{}
	testSummarizeID := "556501c0-a97d-47ed-9add-73b4a4116c83"

	req, err := http.NewRequest("GET", "/api/user/summarize/"+testSummarizeID, bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)

	r := httptest.NewRecorder()

	hc := &Handlers{
		Config: *cfg,
		Repos:  mockStorage,
	}
	hc.GetSummarizeText(r, req)

	assert.Equal(t, http.StatusNotFound, r.Code)
}
