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

	var liked models.Likes
	db.Debug().Unscoped().Where("user_id=? AND post_id=?", userID, postIDUint).Take(&liked)

	if liked.ID != 0 {
		if liked.DeletedAt.Valid {

			if err = db.Unscoped().Model(&liked).Update("deleted_at", nil).Error; err != nil {
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

		} else {
			err = db.Debug().Delete(&liked).Error
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
		}

	} else {

		Likes.CreatedAt = time.Now()
		Likes.UpdatedAt = time.Now()
		Likes.UserID = userID
		Likes.PostID = &postIDUint

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

	var liked models.Likes
	db.Debug().Unscoped().Where("user_id=? AND comment_id=?", userID, commentIDUint).Take(&liked)

	if liked.ID != 0 {
		if liked.DeletedAt.Valid {

			if err = db.Unscoped().Model(&liked).Update("deleted_at", nil).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "BAD_REQUEST",
					"message": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Like post with id %d success", commentID),
			})
			return

		} else {
			err = db.Debug().Delete(&liked).Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "BAD_REQUEST",
					"message": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Unlike post with id %d success", commentID),
			})
			return
		}

	} else {

		Likes.CreatedAt = time.Now()
		Likes.UpdatedAt = time.Now()
		Likes.UserID = userID
		Likes.CommentID = &commentIDUint

		err = db.Debug().Create(&Likes).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "BAD_REQUEST",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Like post with id %d success", commentID),
		})
		return
	}
}
