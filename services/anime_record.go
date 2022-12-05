package services

import (
	"AnimeLifeBackEnd/entity"
	"AnimeLifeBackEnd/entity/response"
	"AnimeLifeBackEnd/global"
	"AnimeLifeBackEnd/wrapper"
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"net/url"
	"strconv"
)

type AnimeRecordService interface {
	FetchAnimeById(bangumiId int) (entity.Anime, error)
	FetchAnimeRecords(userId uint, offset int, limit int) ([]response.AnimeRecord, error)
	FetchAnimeRecordsOfRating(userId uint, offset int, limit int, rating int) ([]response.AnimeRecord, error)
	FetchAnimeRecordSummary(userId uint) (response.AnimeRecordSummary, error)
	AddNewAnime(animeName string) (entity.Anime, error)
	AddNewAnimeRecord(animeId, userId, rating int) (entity.UserAnime, error)
	UpdateAnimeByBangumiId(bangumiId int, anime entity.Anime) (entity.Anime, error)
	UpdateAnimeRecord(animeId, userId, rating int) (entity.UserAnime, error)
	SearchAnimeRecords(userId, offset, limit int, keyword string) ([]response.AnimeRecord, error)
	DeleteAnimeRecord(animeId, userId int) error
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
		Select("animes.*, user_animes.rating AS rating, user_animes.watch_count AS watch_count, user_animes.created_at AS record_at, user_animes.updated_at AS modify_at").
		Offset(offset).
		Limit(limit).
		Order("modify_at DESC, record_at DESC").
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
		Select("animes.*, user_animes.rating AS rating, user_animes.watch_count AS watch_count, user_animes.created_at AS record_at, user_animes.updated_at AS modify_at").
		Offset(offset).
		Limit(limit).
		Order("modify_at DESC, record_at DESC").
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

func (s animeRecordService) AddNewAnime(animeName string) (entity.Anime, error) {
	anime := entity.Anime{
		Name:      animeName,
		NameJp:    "",
		Cover:     "",
		BangumiId: -1,
	}
	// search whether there is a same anime in database first
	var animeInDB entity.Anime
	err := global.MysqlDB.Where("name = ?", animeName).First(&animeInDB).Error
	if err == nil {
		return animeInDB, nil
	}
	if err != gorm.ErrRecordNotFound {
		global.Logger.Errorf("AnimeLifeBackEnd/services/anime_record.go: Unknown error. AddNewAnime: %v", err)
		return anime, err
	}

	encodedAnimeName := url.QueryEscape(animeName)
	resp, err := wrapper.Get("https://api.bgm.tv/search/subject/" + encodedAnimeName + "?type=2&responseGroup=small")
	if err != nil {
		global.Logger.Errorf("AnimeLifeBackEnd/services/anime_record.go: AddNewAnime: %v", err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		var data entity.BangumiResponse
		if err := json.Unmarshal(body, &data); err != nil {
			global.Logger.Errorf("AnimeLifeBackEnd/services/anime_record.go: jsonify when AddNewAnime: %v", err)
		} else {
			animeInfo := data.List[0]
			anime.NameJp = animeInfo.Name
			anime.Cover = animeInfo.Images.Large
			anime.BangumiId = animeInfo.Id
		}
	}

	// add new anime to database
	err = global.MysqlDB.Create(&anime).Error
	if err != nil {
		global.Logger.Errorf("AnimeLifeBackEnd/services/anime_record.go: Fail to add new anime. AddNewAnime: %v", err)
	}

	return anime, err
}

func (s animeRecordService) AddNewAnimeRecord(animeId int, userId int, rating int) (entity.UserAnime, error) {
	userAnime := entity.UserAnime{
		UserId:  userId,
		AnimeId: animeId,
		Rating:  rating,
	}

	// check whether there is a same record in database first
	var userAnimeInDB entity.UserAnime
	err := global.MysqlDB.Where("user_id = ? AND anime_id = ?", userId, animeId).First(&userAnimeInDB).Error
	// if there is a same record, update the watch count and rating, save previous one to history table
	if err == nil {
		animeRecordBak := entity.AnimeRecord{
			UserId:    userAnimeInDB.UserId,
			AnimeId:   userAnimeInDB.AnimeId,
			Rating:    userAnimeInDB.Rating,
			CreatedAt: userAnimeInDB.CreatedAt,
			UpdatedAt: userAnimeInDB.UpdatedAt,
			DeletedAt: userAnimeInDB.DeletedAt,
			Comment:   userAnimeInDB.Comment,
		}
		err = global.MysqlDB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&animeRecordBak).Error; err != nil {
				global.Logger.Errorf("AnimeLifeBackEnd/services/anime_record.go: Fail to add new anime record backup, rolling back. AddNewAnimeRecord: %v", err)
				return err
			}
			if err := tx.Model(&userAnimeInDB).Clauses(clause.Returning{}).Updates(map[string]interface{}{
				"rating":      rating,
				"watch_count": userAnimeInDB.WatchCount + 1,
			}).Error; err != nil {
				global.Logger.Errorf("AnimeLifeBackEnd/services/anime_record.go: Fail to update anime record, rolling back. AddNewAnimeRecord: %v", err)
				return err
			}
			return nil
		})
		userAnime = userAnimeInDB
	} else {
		if err != gorm.ErrRecordNotFound {
			global.Logger.Errorf("AnimeLifeBackEnd/services/anime_record.go: Unknown error when finding user_anime. AddNewAnimeRecord: %v", err)
		}
		// if there is no same record, create a new one
		err = global.MysqlDB.Create(&userAnime).Error
		if err != nil {
			global.Logger.Errorf("AnimeLifeBackEnd/services/anime_record.go: Fail to add new anime record. AddNewAnimeRecord: %v", err)
		}
	}

	return userAnime, err
}

func (s animeRecordService) FetchAnimeById(animeId int) (entity.Anime, error) {
	anime := entity.Anime{
		Model: gorm.Model{},
	}
	anime.ID = uint(animeId)

	err := global.MysqlDB.First(&anime).Error
	if err != nil {
		global.Logger.Errorf("Fail to find anime. Err: %v", err)
	}
	return anime, err
}

func (s animeRecordService) UpdateAnimeByBangumiId(bangumiId int, anime entity.Anime) (entity.Anime, error) {
	// access to bangumi v0 api must contain a user-agent
	resp, err := wrapper.GetWithHeader("https://api.bgm.tv/v0/subjects/"+strconv.Itoa(bangumiId), map[string]string{"User-Agent": "Kmakise/AnimeLife"})
	if err != nil {
		global.Logger.Errorf("Fail to get bangumi info. Err: %v", err)
		return anime, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var data entity.BangumiAnime
	if err := json.Unmarshal(body, &data); err != nil {
		global.Logger.Errorf("Fail to jsonify bangumi info. Err: %v", err)
		return anime, err
	}
	//global.Logger.Debugf("Data: %v", data)
	anime.NameJp = data.Name
	anime.Cover = data.Images.Large
	anime.BangumiId = data.Id
	//global.Logger.Debugf("Anime: %v", anime)
	err = global.MysqlDB.Save(&anime).Error
	if err != nil {
		global.Logger.Errorf("Fail to update anime in DB. Err: %v", err)
	}
	return anime, err
}

func (s animeRecordService) UpdateAnimeRecord(animeId, userId, rating int) (entity.UserAnime, error) {
	userAnime := entity.UserAnime{
		UserId:  userId,
		AnimeId: animeId,
		Rating:  rating,
	}
	// this update doesn't change updated_at
	// using .Clauses(clause.Returning{}) only return changed column, not the whole record
	err := global.MysqlDB.Model(&userAnime).Clauses(clause.Returning{}).UpdateColumn("rating", rating).Error
	if err != nil {
		global.Logger.Errorf("Fail to update rating column of anime record. Err: %v", err)
	}
	global.MysqlDB.First(&userAnime)
	return userAnime, err
}

func (s animeRecordService) SearchAnimeRecords(userId, offset, limit int, keyword string) ([]response.AnimeRecord, error) {
	user := entity.User{
		Model: gorm.Model{},
	}
	user.ID = uint(userId)
	animeRecords := make([]response.AnimeRecord, 0)

	err := global.MysqlDB.Table("animes").
		Joins("JOIN user_animes ON user_animes.anime_id = animes.id").
		Where("user_animes.user_id = ?", userId).
		Where("animes.name LIKE ?", "%"+keyword+"%").
		Or("animes.name_jp LIKE ?", "%"+keyword+"%").
		Select("animes.*, user_animes.rating AS rating, user_animes.watch_count AS watch_count, user_animes.created_at AS record_at, user_animes.updated_at AS modify_at").
		Offset(offset).
		Limit(limit).
		Order("modify_at DESC, record_at DESC").
		Find(&animeRecords).Error
	if err != nil {
		global.Logger.Errorf("%v", err)
	}
	return animeRecords, err
}

func (s animeRecordService) DeleteAnimeRecord(animeId, userId int) error {
	record := entity.UserAnime{
		UserId:  userId,
		AnimeId: animeId,
	}
	// delete record with soft delete can still be found in fetch api
	err := global.MysqlDB.Unscoped().Delete(&record).Error
	if err != nil {
		global.Logger.Errorf("Fail to delete anime record. Err: %v", err)
	}
	return err
}
