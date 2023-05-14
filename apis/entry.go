package apis

import "AnimeLifeBackEnd/services"

type ApiGroup struct {
	User           UserApi
	AnimeRecordApi AnimeRecordApi
	AuthApi        AuthApi
}

var (
	userService        = services.ServiceGroupApp.User
	animeRecordService = services.ServiceGroupApp.AnimeRecord
	authService        = services.ServiceGroupApp.Auth
)

var ApiGroupApp = ApiGroup{
	User:           userApi{},
	AnimeRecordApi: animeRecordApi{},
	AuthApi:        authApi{},
}
