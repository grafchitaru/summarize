package models

import "time"

type Summarize struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Text      *string   `json:"text"`
	Result    *string   `json:"result"`
	Status    string    `json:"status"`
	Tokens    uint64    `json:"tokens"`
}

type NewSummarize struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
	Text   string `json:"text"`
	Status string `json:"status"`
	Tokens int    `json:"tokens"`
}
