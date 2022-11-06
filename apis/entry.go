package apis

import "AnimeLifeBackEnd/services"

type ApiGroup struct {
	Login          LoginApi
	User           UserApi
	AnimeRecordApi AnimeRecordApi
}

var (
	loginService       = services.ServiceGroupApp.Login
	userService        = services.ServiceGroupApp.User
	animeRecordService = services.ServiceGroupApp.AnimeRecord
)

var ApiGroupApp = ApiGroup{
	Login:          loginApi{},
	User:           userApi{},
	AnimeRecordApi: animeRecordApi{},
}
