package api

import (
	"AnimeLifeBackEnd/apis"
	"github.com/gin-gonic/gin"
)

type animeRecordRouter struct{}

func (r animeRecordRouter) AddAnimeRecordRoutes(rg *gin.RouterGroup) {
	animeRecord := rg.Group("/anime_record")
	{
		animeRecord.GET("/:userId", apis.ApiGroupApp.AnimeRecordApi.FetchAnimeRecords)
		animeRecord.GET("/:userId/rating/:rating", apis.ApiGroupApp.AnimeRecordApi.FetchAnimeRecordsOfRating)
		animeRecord.GET("/:userId/summary", apis.ApiGroupApp.AnimeRecordApi.FetchAnimeRecordSummary)
		animeRecord.POST("/:userId/addRecord", apis.ApiGroupApp.AnimeRecordApi.AddAnimeRecord)
		animeRecord.POST("/:userId/updateRecord", apis.ApiGroupApp.AnimeRecordApi.UpdateAnimeRecord)
		animeRecord.GET("/:userId/search", apis.ApiGroupApp.AnimeRecordApi.SearchAnimeRecords)
	}
}
