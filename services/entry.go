package services

type ServiceGroup struct {
	User        UserService
	AnimeRecord AnimeRecordService
}

var ServiceGroupApp = ServiceGroup{
	User:        userService{},
	AnimeRecord: animeRecordService{},
}
