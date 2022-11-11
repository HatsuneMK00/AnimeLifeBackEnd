package middlewares

import (
	"AnimeLifeBackEnd/entity"
	"AnimeLifeBackEnd/entity/request"
	"AnimeLifeBackEnd/global"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

func InitJWTAuth() *jwt.GinJWTMiddleware {
	jwtConfig := global.Config.Jwt
	global.Logger.Debugf("jwt config: %+v", jwtConfig)

	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       jwtConfig.Realm,
		Key:         []byte(jwtConfig.SecretKey),
		Timeout:     12 * time.Hour,
		MaxRefresh:  4 * time.Hour,
		IdentityKey: jwtConfig.IdentityKey,
		SendCookie:  false,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*entity.User); ok {
				return jwt.MapClaims{
					jwtConfig.IdentityKey: v.ID,
					"username":            v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login request.Login
			if err := c.ShouldBindJSON(&login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userName := login.Username
			password := login.Password

			user := entity.User{
				Model: gorm.Model{},
			}
			result := global.MysqlDB.Where("username = ?", userName).First(&user)
			if result.Error != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &user, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": "you are not authorized",
			})
		},
		TokenLookup:       "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:     "Bearer",
		TimeFunc:          time.Now,
		SendAuthorization: false,
	})
	if err != nil {
		global.Logger.Errorf("JWT Error:" + err.Error())
	}
	return middleware
}
