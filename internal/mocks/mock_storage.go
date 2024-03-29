package mocks

import (
	"errors"
	"github.com/grafchitaru/summarize/internal/models"
)

type GetUserFunc func(login string) (string, error)
type GetUserPasswordFunc func(login string) (string, error)
type RegistrationFunc func(id string, login string, password string) (string, error)
type GetSummarizeByTextFunc func(text string) (string, error)
type CreateSummarizeFunc func(summarize models.NewSummarize) error
type GetStatusFunc func(user_id string, AiMaxLimitCount int, AiMaxLimitTokens int) (models.Status, error)
type MockStorage struct {
	PingError              error
	GetUserFunc            GetUserFunc
	RegistrationFunc       RegistrationFunc
	GetUserPasswordFunc    GetUserPasswordFunc
	GetSummarizeByTextFunc GetSummarizeByTextFunc
	CreateSummarizeFunc    CreateSummarizeFunc
	GetStatusFunc          GetStatusFunc
	Users                  map[string]string
	IDs                    map[string]string
	Passwords              map[string]string
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		Users:     make(map[string]string),
		IDs:       make(map[string]string),
		Passwords: make(map[string]string),
	}
}

func (ms *MockStorage) Ping() error {
	return ms.PingError
}

func (ms *MockStorage) Close() {
	// Implementation for Close method
}

func (ms *MockStorage) GetUser(login string) (string, error) {
	if ms.GetUserFunc != nil {
		return ms.GetUserFunc(login)
	}
	user, exists := ms.Users[login]
	if !exists {
		return "", errors.New("user not found")
	}
	return user, nil
}

func (ms *MockStorage) GetUserPassword(login string) (string, error) {
	if ms.GetUserPasswordFunc != nil {
		return ms.GetUserPasswordFunc(login)
	}
	password, exists := ms.Passwords[login]
	if !exists {
		return "", errors.New("user not found")
	}
	return password, nil
}

func (ms *MockStorage) SetUserPassword(login string, password string) {
	ms.Passwords[login] = password
}

func (ms *MockStorage) Registration(id string, login string, password string) (string, error) {
	if ms.RegistrationFunc != nil {
		return ms.RegistrationFunc(id, login, password)
	}
	if _, exists := ms.Users[login]; exists {
		return "", errors.New("user already exists")
	}
	ms.Users[login] = password
	ms.IDs[login] = id
	return id, nil
}

func (ms *MockStorage) UpdateSummarizeStatus(id string, status string) error {
	return nil
}

func (ms *MockStorage) UpdateSummarizeResult(id string, status string, result string) error {
	return nil
}

func (ms *MockStorage) GetSummarize(id string, user_id string) (models.Summarize, error) {
	return models.Summarize{}, nil
}

func (ms *MockStorage) GetSummarizeByText(text string) (string, error) {
	if ms.GetSummarizeByTextFunc != nil {
		return ms.GetSummarizeByTextFunc(text)
	}
	return "", errors.New("not found")
}

func (ms *MockStorage) CreateSummarize(summarize models.NewSummarize) error {
	if ms.CreateSummarizeFunc != nil {
		return ms.CreateSummarizeFunc(summarize)
	}
	return nil
}

func (ms *MockStorage) GetStat(user_id string) ([]models.Stat, error) {
	return []models.Stat{}, nil
}

func (ms *MockStorage) GetStatus(user_id string, AiMaxLimitCount int, AiMaxLimitTokens int) (models.Status, error) {
	return models.Status{}, nil
}
