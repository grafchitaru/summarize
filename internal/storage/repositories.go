package storage

import "time"

type Summarize struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
	Result    string    `json:"result"`
	Status    string    `json:"status"`
	Tokens    uint64    `json:"tokens"`
}

type Repositories interface {
	Ping() error
	Close()
	GetUser(login string) (string, error)
	GetUserPassword(login string) (string, error)
	Registration(id string, login string, password string) (string, error)
	CreateSummarize(id string, user_id string, text string, status string, tokens int) error
	UpdateSummarizeStatus(id string, status string) error
	UpdateSummarizeResult(id string, status string, result string) error
	GetSummarize(id string, user_id string) (Summarize, error)
	GetSummarizeByText(text string) (string, error)
}
