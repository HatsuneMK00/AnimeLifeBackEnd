package config

type Server struct {
	Zap    Zap    `json:"zap" yaml:"zap"`
	System System `json:"system" yaml:"system"`
	Mysql  Mysql  `json:"mysql" yaml:"mysql"`
	Jwt    Jwt    `json:"jwt" yaml:"jwt"`
	Redis  Redis  `json:"redis" yaml:"redis"`
	Email  Email  `json:"email" yaml:"email"`
}
