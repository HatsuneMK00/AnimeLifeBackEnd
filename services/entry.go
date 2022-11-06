package services

type ServiceGroup struct {
	User        UserService
	Login       LoginService
	AnimeRecord AnimeRecordService
}

var ServiceGroupApp = ServiceGroup{
	User:        userService{},
	Login:       loginService{},
	AnimeRecord: animeRecordService{},
}
