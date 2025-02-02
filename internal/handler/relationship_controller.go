package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"twittir-go/internal/database"
	"twittir-go/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func FollowUser(c *gin.Context) {
	db := database.GetDB()
	Relationship := domain.Relationship{}

	// Current User ID
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	// User ID who want to follow by current user
	followingIDStr := c.Param("id")
	followingID, err := strconv.ParseUint(followingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid followingID"})
		return
	}

	followingIDUint := uint(followingID)

	Relationship.FollowerID = userID
	Relationship.FollowingID = followingIDUint

	var followed domain.Relationship
	err = db.Debug().Where("follower_id=?", userID).Where("following_id=?", followingIDUint).Take(&followed).Error

	if followed.ID != 0 {
		err = db.Debug().Where("follower_id=?", userID).Where("following_id=?", followingIDUint).Delete(&followed).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "BAD_REQUEST",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Unfollow user with id %d success", followingID),
		})
		return
	} else {
		err = db.Debug().Create(&Relationship).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "BAD_REQUEST",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully followed user with id %d", followingID),
		})
		return
	}
}
