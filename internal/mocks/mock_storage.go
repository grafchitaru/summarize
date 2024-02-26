package mocks

import (
	"errors"
	"github.com/grafchitaru/summarize/internal/storage"
)

type GetUserFunc func(login string) (string, error)
type GetUserPasswordFunc func(login string) (string, error)
type RegistrationFunc func(id string, login string, password string) (string, error)
type GetSummarizeByTextFunc func(text string) (string, error)
type CreateSummarizeFunc func(id string, user_id string, text string, status string, tokens int) error

type MockStorage struct {
	PingError              error
	GetUserFunc            GetUserFunc
	RegistrationFunc       RegistrationFunc
	GetUserPasswordFunc    GetUserPasswordFunc
	GetSummarizeByTextFunc GetSummarizeByTextFunc
	CreateSummarizeFunc    CreateSummarizeFunc
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

func (ms *MockStorage) GetSummarize(id string, user_id string) (storage.Summarize, error) {
	return storage.Summarize{}, nil
}

func (ms *MockStorage) GetSummarizeByText(text string) (string, error) {
	if ms.GetSummarizeByTextFunc != nil {
		return ms.GetSummarizeByTextFunc(text)
	}
	return "", errors.New("not found")
}

func (ms *MockStorage) CreateSummarize(id string, userID string, text string, status string, tokens int) error {
	if ms.CreateSummarizeFunc != nil {
		return ms.CreateSummarizeFunc(id, userID, text, status, tokens)
	}
	return nil
}

func (ms *MockStorage) GetStat(user_id string) ([]storage.Stat, error) {
	return []storage.Stat{}, nil
}
