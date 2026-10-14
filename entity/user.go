package entity

import "github.com/google/uuid"

type User struct {
	Id         uuid.UUID
	Email      string
	CreatedAt  string `json:"created_at"`
	LastReadId int    `json:"last_read_id"`
	Statistic  UserStatistic
}

type UserStatistic struct {
	Total   int
	ByMonth map[string]int `json:"by_month"`
	ByDays  map[string]int `json:"by_days"`
}
