package core

import (
	"AnimeLifeBackEnd/entity"
	"AnimeLifeBackEnd/global"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitMysqlDB() *gorm.DB {
	c := global.Config.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.Username, c.Password, c.Path, c.Port, c.Dbname, c.Config)
	global.Logger.Infof("%s", dsn)

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err != nil {
		global.Logger.Errorf("%s", err)
		return nil
	}
	return db
}

func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(&entity.User{}, &entity.Anime{}, &entity.UserAnime{}, &entity.AnimeRecord{})
	if err != nil {
		global.Logger.Errorf("Database: RegisterTables failed, err: %v", zap.Error(err))
		os.Exit(0)
	}
	err = db.SetupJoinTable(&entity.User{}, "Animes", &entity.UserAnime{})
	if err != nil {
		global.Logger.Errorf("Database: SetupJoinTable failed, err: %v", zap.Error(err))
		os.Exit(0)
	}
	err = db.SetupJoinTable(&entity.Anime{}, "Users", &entity.UserAnime{})
	if err != nil {
		global.Logger.Errorf("Database: SetupJoinTable failed, err: %v", zap.Error(err))
		os.Exit(0)
	}
	global.Logger.Info("Database: Register table successfully")
}
