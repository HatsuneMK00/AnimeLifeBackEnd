package entity

import (
	"gorm.io/gorm"
	"time"
)

type UserAnime struct {
	AnimeId   int            `json:"anime_id" gorm:"primary_key"`
	UserId    int            `json:"user_id" gorm:"primary_key"`
	Rating    int            `json:"rating"`
	Comment   string         `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
