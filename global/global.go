package global

import (
	"AnimeLifeBackEnd/config"
	"AnimeLifeBackEnd/websocket/base"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger  *zap.SugaredLogger
	Config  config.Server
	Viper   *viper.Viper
	MysqlDB *gorm.DB
	WsHub   base.Hub
)
