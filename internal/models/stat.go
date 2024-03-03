package models

type Stat struct {
	UserId string `json:"user_id"`
	Status string `json:"status"`
	Count  uint64 `json:"count"`
	Tokens uint64 `json:"tokens"`
}
