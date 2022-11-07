package services

import (
	"AnimeLifeBackEnd/entity"
	"AnimeLifeBackEnd/entity/response"
	"AnimeLifeBackEnd/global"
	"gorm.io/gorm"
)

type AnimeRecordService interface {
	FetchAnimeRecords(userId uint, offset int, limit int) ([]response.AnimeRecord, error)
	FetchAnimeRecordsOfRating(userId uint, offset int, limit int, rating int) ([]response.AnimeRecord, error)
	FetchAnimeRecordSummary(userId uint) (response.AnimeRecordSummary, error)
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

func (s animeRecordService) FetchAnimeRecordsOfRating(userId uint, offset int, limit int, rating int) ([]response.AnimeRecord, error) {
	user := entity.User{
		Model: gorm.Model{},
	}
	user.ID = userId
	animeRecords := make([]response.AnimeRecord, 0)

	err := global.MysqlDB.Table("animes").
		Joins("JOIN user_animes ON user_animes.anime_id = animes.id").
		Where("user_animes.user_id = ?", userId).
		Where("user_animes.rating = ?", rating).
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

func (s animeRecordService) FetchAnimeRecordSummary(userId uint) (response.AnimeRecordSummary, error) {
	user := entity.User{
		Model: gorm.Model{},
	}
	user.ID = userId
	animeRecordSummary := response.AnimeRecordSummary{}

	err := global.MysqlDB.Table("user_animes").
		Where("user_animes.user_id = ?", userId).
		Select("SUM(rating=1) AS rating_one_count,SUM(rating=2) AS rating_two_count," +
			"SUM(rating=3) AS rating_three_count,SUM(rating=4) AS rating_four_count,COUNT(rating) AS total_count").
		Find(&animeRecordSummary).Error
	if err != nil {
		global.Logger.Errorf("%v", err)
	}
	return animeRecordSummary, err
}
