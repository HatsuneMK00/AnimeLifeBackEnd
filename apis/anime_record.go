package apis

import (
	"AnimeLifeBackEnd/entity"
	"AnimeLifeBackEnd/entity/request"
	"AnimeLifeBackEnd/global"
	entity2 "AnimeLifeBackEnd/websocket/base"
	"errors"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/url"
	"strconv"
)

type AnimeRecordApi interface {
	FetchAnimeRecords(c *gin.Context)
	FetchAnimeRecordsOfRating(c *gin.Context)
	FetchAnimeRecordSummary(c *gin.Context)
	FetchHistoryRatingOfRecord(c *gin.Context)
	AddAnimeRecord(c *gin.Context)
	UpdateAnimeRecord(c *gin.Context)
	SearchAnimeRecords(c *gin.Context)
	DeleteAnimeRecord(c *gin.Context)
}

type animeRecordApi struct{}

func (a animeRecordApi) FetchAnimeRecords(c *gin.Context) {
	// return 15 records from offset at a time
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
		})
		global.Logger.Errorf("fail to get user id from jwt token, or user id is not int: %v", err)
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "offset needs to be uint",
		})
		global.Logger.Errorf("offset needs to be uint: %v", err)
		return
	}
	animes, err := animeRecordService.FetchAnimeRecords(uint(userId), offset, 15)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "db error",
		})
		global.Logger.Errorf("db error: %v", err)
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    animes,
		"hasMore": !(len(animes) < 15),
	})
}

func (a animeRecordApi) FetchAnimeRecordsOfRating(c *gin.Context) {
	// return 15 records from offset at a time
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
		})
		global.Logger.Errorf("fail to get user id from jwt token, or user id is not int: %v", err)
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "offset needs to be uint",
		})
		global.Logger.Errorf("offset needs to be uint: %v", err)
		return
	}
	rating, err := strconv.Atoi(c.Param("rating"))
	if err != nil || rating < -1 || rating > 4 {
		c.JSON(400, gin.H{
			"message": "rating needs to be uint and between [-1, 4]",
		})
		global.Logger.Errorf("rating needs to be uint and between [-1, 4]: %v", err)
		return
	}
	animes, err := animeRecordService.FetchAnimeRecordsOfRating(uint(userId), offset, 15, rating)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "db error",
		})
		global.Logger.Errorf("db error: %v", err)
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    animes,
		"hasMore": !(len(animes) < 15),
	})
}

func (a animeRecordApi) FetchAnimeRecordSummary(c *gin.Context) {
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
		})
		return
	}
	summary, err := animeRecordService.FetchAnimeRecordSummary(uint(userId))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "db error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    summary,
	})
}

func (a animeRecordApi) FetchHistoryRatingOfRecord(c *gin.Context) {
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
		})
		return
	}
	animeId, err := strconv.Atoi(c.Query("animeId"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "animeId needs to be uint",
		})
		return
	}
	ratings, err := animeRecordService.FetchHistoryRatingOfRecord(uint(userId), uint(animeId))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "db error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    ratings,
	})
}

func (a animeRecordApi) AddAnimeRecord(c *gin.Context) {
	//	read request body from json
	var record request.AnimeRecordRequest
	err := c.ShouldBindJSON(&record)
	if err != nil {
		global.Logger.Error(err)
		c.JSON(400, gin.H{
			"message": "invalid request body",
		})
		return
	}
	//	validate request body
	if !(record.AnimeRating == -1 || // 想看
		record.AnimeRating == 1 || // 一般
		record.AnimeRating == 2 || // 不错
		record.AnimeRating == 3 || // 好看
		record.AnimeRating == 4) { // 神作
		c.JSON(400, gin.H{
			"message": "invalid anime rating",
		})
		return
	}
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
		})
		return
	}

	global.WsHub.Comm() <- &entity2.Message{Type: "message", Data: "fetching anime info from bangumi api"}
	anime, err := animeRecordService.AddNewAnime(record.AnimeName)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "fail to add anime",
		})
		return
	}
	global.WsHub.Comm() <- &entity2.Message{Type: "message", Data: "anime info fetched"}
	animeRecord, isNewRecord, err := animeRecordService.AddNewAnimeRecord(int(anime.ID), userId, record.AnimeRating, record.Commment)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "fail to add anime record, some error occurred",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data": map[string]interface{}{
			"anime":  anime,
			"record": animeRecord,
			"is_new": isNewRecord,
		},
	})
}

func (a animeRecordApi) UpdateAnimeRecord(c *gin.Context) {
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
		})
		return
	}
	updateRequest := request.AnimeRecordUpdateRequest{}
	err = c.ShouldBindJSON(&updateRequest)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid request body",
		})
		return
	}

	// update anime info if bangumiId is changed
	anime, err := animeRecordService.FetchAnimeById(updateRequest.AnimeId)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "fail to fetch anime",
		})
		return
	}
	shouldAnimeUpdate := false
	isAnimeUpdateFailed := false
	var updatedAnime entity.Anime
	if anime.BangumiId != updateRequest.BangumiId {
		shouldAnimeUpdate = true
		updatedAnime, err = animeRecordService.UpdateAnimeByBangumiId(updateRequest.BangumiId, anime)
		if err != nil {
			isAnimeUpdateFailed = true
		}
	}
	record, err := animeRecordService.UpdateAnimeRecord(updateRequest.AnimeId, userId, updateRequest.AnimeRating, updateRequest.Comment)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "fail to update anime record",
		})
		return
	}
	if !shouldAnimeUpdate || isAnimeUpdateFailed {
		updatedAnime = anime
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data": map[string]interface{}{
			"record":                 record,
			"updated_anime":          updatedAnime,
			"is_anime_update_failed": isAnimeUpdateFailed,
		},
	})
}

func (a animeRecordApi) SearchAnimeRecords(c *gin.Context) {
	// return 15 records from offset at a time
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
		})
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "offset needs to be uint",
		})
		return
	}
	searchText, err := url.QueryUnescape(c.DefaultQuery("searchText", ""))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "searchText needs to be url encoded",
		})
		return
	}
	animeRecords, err := animeRecordService.SearchAnimeRecords(userId, offset, 15, searchText)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "db error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    animeRecords,
		"hasMore": !(len(animeRecords) < 15),
	})
}

func (a animeRecordApi) DeleteAnimeRecord(c *gin.Context) {
	userId, err := getUserIdFromJwtToken(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "fail to get user id from jwt token, or user id is not int",
		})
		return
	}
	deleteRequest := request.AnimeRecordDeleteRequest{}
	err = c.ShouldBindJSON(&deleteRequest)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid request body",
		})
		return
	}
	err = animeRecordService.DeleteAnimeRecord(deleteRequest.AnimeId, userId)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "fail to delete anime record",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success deleting anime record with anime id " + strconv.Itoa(deleteRequest.AnimeId),
	})
}

func getUserIdFromJwtToken(c *gin.Context) (int, error) {
	claims := jwt.ExtractClaims(c)
	userId := claims[global.Config.Jwt.IdentityKey]
	if userId, ok := userId.(float64); ok {
		return int(userId), nil
	} else {
		err := errors.New("user id is not int")
		return 0, err
	}
}
