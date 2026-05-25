package models

import "time"

type YvUser struct {
	Id        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	GoogleId  string    `json:"google_id" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"not null"`
	AvatarUrl *string   `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (YvUser) TableName() string { return "yv_user" }
