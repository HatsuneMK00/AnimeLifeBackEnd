package config

type Jwt struct {
	IdentityKey string `json:"identity-key" yaml:"identity-key" mapstructure:"identity-key"` // 用户标识
	SecretKey   string `json:"secret-key" yaml:"secret-key" mapstructure:"secret-key"`       // 加密token密钥
	Realm       string `json:"realm" yaml:"realm"`                                           // App名称
}
