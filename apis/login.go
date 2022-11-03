package apis

type LoginApi interface {
	LogoutUser(userId string)
}

type loginApi struct{}

func (c loginApi) LogoutUser(userId string) {
	//TODO implement me
	panic("implement me")
}
