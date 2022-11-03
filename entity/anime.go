package entity

import "gorm.io/gorm"

type Anime struct {
	gorm.Model
	Name   string `json:"name" binding:"required"`
	NameJp string `json:"name_jp"`
	Cover  string `json:"cover"`
}
