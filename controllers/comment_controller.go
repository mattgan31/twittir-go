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

func CreateComment(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "BAD_REQUEST",
			"error":  "Invalid PostID",
		})
		return
	}

	postIDUint := uint(postID)
	Comment := models.Comment{}

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.CreatedAt = time.Now()
	Comment.UpdatedAt = time.Now()
	Comment.UserID = userID
	Comment.PostID = postIDUint

	err = db.Debug().Create(&Comment).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comment": gin.H{
			"id":        Comment.ID,
			"comment":   Comment.Description,
			"post_id":   Comment.PostID,
			"createdAt": Comment.CreatedAt,
			"user":      Comment.UserID,
		},
	})
}

func DeleteComment(c *gin.Context) {
	db := database.GetDB()

	var comment models.Comment

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "BAD_REQUEST",
			"error":  "Invalid CommentID",
		})
		return
	}

	commentIDUint := uint(commentID)

	if err := db.Debug().Where("id=?", commentIDUint).
		Where("user_id=?", userID).
		Take(&comment).
		Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	if err := db.Debug().
		Delete(&comment).
		Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Comment with id %d deleted successfully", commentIDUint),
	})
}
