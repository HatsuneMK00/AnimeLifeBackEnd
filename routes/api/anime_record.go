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
	}
}
