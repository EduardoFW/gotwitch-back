package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/lib/pq"
)

type Stream struct {
	Id           string   `json:"id" "gorm:"primary_key"`
	UserId       string   `json:"user_id"`
	UserLogin    string   `json:"user_login"`
	UserName     string   `json:"user_name"`
	GameId       string   `json:"game_id"`
	GameName     string   `json:"game_name"`
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	ViewerCount  int      `json:"viewer_count"`
	StartedAt    string   `json:"started_at"`
	Language     string   `json:"language"`
	ThumbnailUrl string   `json:"thumbnail_url"`
	TagIds       pq.StringArray `json:"tag_ids" gorm:"type:text[]"`
	IsMature     bool     `json:"is_mature"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}