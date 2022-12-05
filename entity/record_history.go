package entity

import (
	"gorm.io/gorm"
	"time"
)

type AnimeRecord struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	AnimeId   int            `json:"anime_id" gorm:"index"`
	UserId    int            `json:"user_id" gorm:"index"`
	Rating    int            `json:"rating"`
	Comment   string         `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
