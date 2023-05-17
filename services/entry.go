package services

type ServiceGroup struct {
	User        UserService
	AnimeRecord AnimeRecordService
	Auth        AuthService
}

var ServiceGroupApp = ServiceGroup{
	User:        userService{},
	AnimeRecord: animeRecordService{},
	Auth:        authService{},
}
