package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string   `json:"username" binding:"required"`
	Password string   `json:"password" binding:"required"`
	Email    string   `json:"email"`
	Animes   []*Anime `json:"animes" gorm:"many2many:user_animes;"`
	ClerkId  string   `json:"clerk_id"`
}
