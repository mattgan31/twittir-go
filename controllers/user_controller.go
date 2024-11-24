package controllers

import (
	"net/http"
	"strconv"
	"twittir-go/database"
	"twittir-go/helpers"
	"twittir-go/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	appJSON = "application/json"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUser struct {
	Username  string `json:"username"`
	Full_Name string `json:"fullname"`
	Bio       string `json:"bio"`
}

type formatUserRegister struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Full_Name     string `json:"fullname"`
	Password      string `json:"password"`
	PasswordVerif string `json:"password_verif"`
}

func UserRegister(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	formatUser := formatUserRegister{}

	if contentType == appJSON {
		c.ShouldBindJSON(&formatUser)
	} else {
		c.ShouldBind(&formatUser)
	}

	if formatUser.Password != formatUser.PasswordVerif {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid password",
		})
		return
	}

	User := models.User{
		Username:  formatUser.Username,
		Email:     formatUser.Email,
		Full_Name: formatUser.Full_Name,
		Password:  formatUser.Password,
	}

	err := db.Debug().Create(&User).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data": gin.H{
			"id":       User.ID,
			"fullname": User.Full_Name,
			"email":    User.Email,
		},
	})
}

func UserLogin(c *gin.Context) {
	var Login Login
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	User := models.User{}
	password := ""
	if contentType == appJSON {
		c.ShouldBindJSON(&Login)
	} else {
		c.ShouldBind(&Login)
	}

	password = Login.Password

	err := db.Debug().Where("username=?", Login.Username).Take(&User).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "UNAUTHORIZED",
			"message": "Invalid username/password",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))

	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "UNAUTHORIZED",
			"message": "Invalid username/password",
		})
		return
	}

	token := helpers.GenerateToken(User.ID, User.Username)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func SettingsProfile(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))
	User := models.User{}

	if err := db.First(&User, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	updateUserProfile := models.User{}

	if contentType == appJSON {
		c.ShouldBindJSON(&updateUserProfile)
	} else {
		c.ShouldBind(&updateUserProfile)
	}

	if err := db.Model(&User).Updates(UpdateUser{Full_Name: updateUserProfile.Full_Name, Username: updateUserProfile.Username, Bio: updateUserProfile.Bio}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"username":  User.Username,
			"fullname": User.Full_Name,
			"bio":       User.Bio,
		},
	})
}

func GetDetailUser(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))
	User := models.User{}

	err := db.First(&User, userID).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              User.ID,
		"fullname":        User.Full_Name,
		"username":        User.Username,
		"profile_picture": User.Profile_Picture,
	})
}

func SearchUser(c *gin.Context) {
	db := database.GetDB()

	usernameParam := c.DefaultQuery("username", "")

	var User []models.User

	err := db.Debug().Where("username like ?", usernameParam+"%").Find(&User).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	response := make([]FormatUsers, len(User))

	for i, user := range User {
		formattedUsers := FormatUsers{
			ID:       user.ID,
			Username: user.Username,
		}

		response[i] = formattedUsers
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func GetUserByID(c *gin.Context) {
	db := database.GetDB()

	var user models.User

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "BAD_REQUEST",
			"error":  "Invalid UserID",
		})
		return
	}

	userIDUint := uint(userID)

	if err := db.Debug().Where("id=?", userIDUint).Take(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":              user.ID,
			"fullname":        user.Full_Name,
			"username":        user.Username,
			"profile_picture": user.Profile_Picture,
		},
	})
}
