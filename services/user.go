package services

import (
	"AnimeLifeBackEnd/entity"
	"AnimeLifeBackEnd/global"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	FindUser(id uint) (*entity.User, bool)
	FindUserByEmail(email string) (*entity.User, bool)
	FindUsersWithOffset(offset int) ([]entity.User, bool)
	AddUser(user *entity.User) (*entity.User, int64)
}

type userService struct{}

func (s userService) AddUser(user *entity.User) (*entity.User, int64) {
	// encrypt password with bcrypt
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		global.Logger.Errorf("%v", err)
		return user, 0
	}
	user.Password = string(hashedPwd)

	result := global.MysqlDB.Create(user)
	if result.Error != nil {
		global.Logger.Errorf("%v", result.Error)
	}
	user.Password = ""
	return user, result.RowsAffected
}

func (s userService) FindUsersWithOffset(offset int) ([]entity.User, bool) {
	users := make([]entity.User, 0)
	result := global.MysqlDB.Order("id desc").Limit(10).Offset(offset).Find(&users)
	ok := true
	if result.Error != nil {
		global.Logger.Errorf("%v", result.Error)
		ok = false
	}
	for i := range users {
		users[i].Password = ""
	}
	return users, ok
}

func (s userService) FindUser(id uint) (*entity.User, bool) {
	user := entity.User{
		Model: gorm.Model{},
	}
	result := global.MysqlDB.First(&user, id)
	ok := true
	if result.Error != nil {
		global.Logger.Errorf("%v", result.Error)
		ok = false
	}
	user.Password = ""
	return &user, ok
}

func (s userService) FindUserByEmail(email string) (*entity.User, bool) {
	user := entity.User{
		Model: gorm.Model{},
	}
	result := global.MysqlDB.Where("email = ?", email).First(&user)
	ok := true
	if result.Error != nil {
		global.Logger.Errorf("%v", result.Error)
		ok = false
	}
	user.Password = ""
	return &user, ok
}
