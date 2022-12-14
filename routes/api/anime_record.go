package api

import (
	"AnimeLifeBackEnd/apis"
	"github.com/gin-gonic/gin"
)

type animeRecordRouter struct{}

func (r animeRecordRouter) AddAnimeRecordRoutes(rg *gin.RouterGroup) {
	animeRecord := rg.Group("/anime_record")
	{
		animeRecord.GET("", apis.ApiGroupApp.AnimeRecordApi.FetchAnimeRecords)
		animeRecord.GET("/rating/:rating", apis.ApiGroupApp.AnimeRecordApi.FetchAnimeRecordsOfRating)
		animeRecord.GET("/summary", apis.ApiGroupApp.AnimeRecordApi.FetchAnimeRecordSummary)
		animeRecord.GET("/history_rating", apis.ApiGroupApp.AnimeRecordApi.FetchHistoryRatingOfRecord)
		animeRecord.POST("/addRecord", apis.ApiGroupApp.AnimeRecordApi.AddAnimeRecord)
		animeRecord.POST("/updateRecord", apis.ApiGroupApp.AnimeRecordApi.UpdateAnimeRecord)
		animeRecord.GET("/search", apis.ApiGroupApp.AnimeRecordApi.SearchAnimeRecords)
		animeRecord.POST("/deleteRecord", apis.ApiGroupApp.AnimeRecordApi.DeleteAnimeRecord)
	}
}
