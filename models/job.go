package models

import (
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
