package storage

import "github.com/grafchitaru/summarize/internal/models"

type Repositories interface {
	Ping() error
	Close()
	GetUser(login string) (string, error)
	GetUserPassword(login string) (string, error)
	Registration(id string, login string, password string) (string, error)
	CreateSummarize(id string, user_id string, text string, status string, tokens int) error
	UpdateSummarizeStatus(id string, status string) error
	UpdateSummarizeResult(id string, status string, result string) error
	GetSummarize(id string, user_id string) (models.Summarize, error)
	GetSummarizeByText(text string) (string, error)
	GetStat(user_id string) ([]models.Stat, error)
	GetStatus(user_id string, AiMaxLimitCount int, AiMaxLimitTokens int) (models.Status, error)
}
