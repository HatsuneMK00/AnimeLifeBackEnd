package apis

import (
	"AnimeLifeBackEnd/entity"
	"AnimeLifeBackEnd/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"net/smtp"
	"strconv"
)

type AuthApi interface {
	LoginViaEmail(c *gin.Context)
	RegisterViaEmail(c *gin.Context)
}

type authApi struct{}

func (a authApi) LoginViaEmail(c *gin.Context) {
	// 1. Get email from request body
	email := c.Query("email")
	global.Logger.Infof("email: %v", email)
	// 2. Check if email exists in database
	if _, ok := userService.FindUserByEmail(email); ok {
		// 3. generate a random 6 digits code and send it to the email
		// 3.1 generate a random 6 digits code
		code := rand.Intn(900000) + 100000
		msgBody := "Your verification code for AnimeLife is %v. This code will expire in 5 minutes."
		msgBody = fmt.Sprintf(msgBody, code)
		msg := []byte("To: " + email + "\r\n" +
			"Subject: AnimeLife Verification Code\r\n" +
			"\r\n" +
			msgBody + "\r\n")
		// 3.2 send the code to the email
		auth := smtp.PlainAuth("", global.Config.Email.Username, global.Config.Email.Password, global.Config.Email.Host)
		err := smtp.SendMail(global.Config.Email.Host+":"+strconv.Itoa(global.Config.Email.Port), auth, global.Config.Email.Username, []string{email}, msg)
		if err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"message": "failed to send email",
			})
			global.Logger.Errorf("failed to send email: %v", err)
			return
		}
		// 4. save the code to redis
		err = authService.SaveCode(email, code)
		if err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"message": "failed to save code to redis",
			})
			global.Logger.Errorf("failed to save code to redis: %v", err)
			return
		}
		// 5. return success message
		c.JSON(200, gin.H{
			"code":    200,
			"message": "successfully send verification code to email",
		})
	} else {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "user not found",
		})
	}
}

func (a authApi) RegisterViaEmail(c *gin.Context) {
	username := c.Query("username")
	email := c.Query("email")
	password := c.Query("password")

	toAdd := entity.User{
		Model:    gorm.Model{},
		Username: username,
		Email:    email,
		Password: password,
	}
	// Check if email exists in database
	if _, ok := userService.FindUserByEmail(email); ok {
		c.JSON(409, gin.H{
			"code":    409,
			"message": "email already exists",
		})
		return
	}
	user, rowAffected := userService.AddUser(&toAdd)
	if rowAffected == 1 {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "successfully register",
		})
		global.Logger.Infof("successfully register: %v", user)
	} else {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "failed to register",
		})
		global.Logger.Errorf("failed to register: %v", user)
	}
}
