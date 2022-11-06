package apis

import "AnimeLifeBackEnd/services"

type ApiGroup struct {
	User           UserApi
	AnimeRecordApi AnimeRecordApi
}

var (
	userService        = services.ServiceGroupApp.User
	animeRecordService = services.ServiceGroupApp.AnimeRecord
)

var ApiGroupApp = ApiGroup{
	User:           userApi{},
	AnimeRecordApi: animeRecordApi{},
}
