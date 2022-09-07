package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	Id        string `json:"id" "gorm:"primary_key"`
	Name      string `json:"name"`
	BoxArtUrl string `json:"box_art_url"`
	JobID     int    `json:"-"`
	Job       Job    `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
