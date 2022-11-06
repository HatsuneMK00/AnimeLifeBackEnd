package entity

import "gorm.io/gorm"

type Anime struct {
	gorm.Model
	Name      string  `json:"name" binding:"required"`
	NameJp    string  `json:"name_jp"`
	Cover     string  `json:"cover"`
	BangumiId int     `json:"bangumi_id"`
	Users     []*User `json:"users" gorm:"many2many:user_animes;"`
}
