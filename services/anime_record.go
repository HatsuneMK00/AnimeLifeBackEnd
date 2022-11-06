package services

import (
	"AnimeLifeBackEnd/entity"
	"AnimeLifeBackEnd/entity/response"
	"AnimeLifeBackEnd/global"
	"gorm.io/gorm"
)

type AnimeRecordService interface {
	FetchAnimeRecords(userId uint, offset int, limit int) ([]response.AnimeRecord, error)
}

type animeRecordService struct{}

func (s animeRecordService) FetchAnimeRecords(userId uint, offset int, limit int) ([]response.AnimeRecord, error) {
	user := entity.User{
		Model: gorm.Model{},
	}
	user.ID = userId
	animeRecords := make([]response.AnimeRecord, 0)

	// use log to print the sql for debug
	//tx := global.MysqlDB.Session(&gorm.Session{Logger: logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags),
	//	logger.Config{
	//		SlowThreshold:             200 * time.Millisecond,
	//		LogLevel:                  logger.Info,
	//		Colorful:                  true,
	//		IgnoreRecordNotFoundError: true,
	//	},
	//)})
	//err := tx.Table("animes").
	//	Joins("JOIN user_animes ON user_animes.anime_id = animes.id").
	//	Where("user_animes.user_id = ?", userId).
	//	Select("animes.*, user_animes.rating, user_animes.created_at AS record_at").
	//	Offset(offset).
	//	Limit(limit).
	//	Order("record_at DESC").
	//	Find(&animes).Error

	err := global.MysqlDB.Table("animes").
		Joins("JOIN user_animes ON user_animes.anime_id = animes.id").
		Where("user_animes.user_id = ?", userId).
		Select("animes.*, user_animes.rating AS rating, user_animes.created_at AS record_at").
		Offset(offset).
		Limit(limit).
		Order("record_at DESC").
		Find(&animeRecords).Error
	if err != nil {
		global.Logger.Errorf("%v", err)
	}
	return animeRecords, err
}
