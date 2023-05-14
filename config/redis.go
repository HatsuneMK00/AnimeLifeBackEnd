package config

type Redis struct {
	Addr     string `json:"addr" yaml:"addr"`         // 服务器地址
	Password string `json:"password" yaml:"password"` // 密码
	DB       int    `json:"db" yaml:"db"`             // 数据库
}
