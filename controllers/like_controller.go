package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"twittir-go/database"
	"twittir-go/helpers"
	"twittir-go/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CreateLikePost(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	Likes := models.Likes{}

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PostID"})
		return
	}

	postIDUint := uint(postID)

	if contentType == appJSON {
		c.ShouldBindJSON(&Likes)
	} else {
		c.ShouldBind(&Likes)
	}

	Likes.Created_At = time.Now()
	Likes.Updated_At = time.Now()
	Likes.UserID = userID
	Likes.PostID = &postIDUint

	var liked models.Likes
	db.Debug().Where("user_id=?", Likes.UserID).Where("post_id=?", Likes.PostID).Take(&liked)

	if liked.ID != 0 {
		err = db.Debug().Where("user_id=?", Likes.UserID).Where("post_id=?", Likes.PostID).Delete(&liked).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Unlike post with id %d success", postID),
		})
		return
	} else {
		err = db.Debug().Create(&Likes).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Like post with id %d success", postID),
		})
		return
	}
}

func CreateLikeComment(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	Likes := models.Likes{}

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid commentID"})
		return
	}

	commentIDUint := uint(commentID)

	if contentType == appJSON {
		c.ShouldBindJSON(&Likes)
	} else {
		c.ShouldBind(&Likes)
	}

	Likes.Created_At = time.Now()
	Likes.Updated_At = time.Now()
	Likes.UserID = userID
	Likes.CommentID = &commentIDUint

	var liked models.Likes
	db.Debug().Where("user_id=?", Likes.UserID).Where("comment_id=?", Likes.CommentID).Take(&liked)

	if liked.ID != 0 {
		err = db.Debug().Where("user_id=?", Likes.UserID).Where("comment_id=?", Likes.CommentID).Delete(&liked).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Unlike comment with id %d success", commentID),
		})
		return
	} else {
		err = db.Debug().Create(&Likes).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Like comment with id %d success", commentID),
		})
		return
	}
}
