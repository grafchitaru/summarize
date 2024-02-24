package mocks

import (
	"errors"
)

type GetUserFunc func(login string) (string, error)
type RegistrationFunc func(id string, login string, password string) (string, error)

type MockStorage struct {
	PingError        error
	GetUserFunc      GetUserFunc
	RegistrationFunc RegistrationFunc
	Users            map[string]string
	IDs              map[string]string
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		Users: make(map[string]string),
		IDs:   make(map[string]string),
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
