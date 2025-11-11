package models

import "time"

type Log struct {
	Id          *int64     `json:"id"`
	UserId      *int64     `json:"id_user"`
	Username    *string    `json:"username"`
	Action      *string    `json:"action"`
	Description *string    `json:"description"`
	Url         *string    `json:"url"`
	CreatedAt   *time.Time `json:"created_at"`
}
