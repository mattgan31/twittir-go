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
	"gorm.io/gorm"
)

// Type for Formatting
type FormatPosts struct {
	ID        uint             `json:"id"`
	Post      string           `json:"post"`
	CreatedAt time.Time        `json:"createdAt"`
	User      FormatUsers      `json:"user"`
	Likes     []FormatLikes    `json:"likes"`
	Comment   []FormatComments `json:"comments"`
}

type FormatComments struct {
	ID          uint          `json:"id"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"createdAt"`
	User        FormatUsers   `json:"user"`
	Likes       []FormatLikes `json:"likes"`
}

type FormatUsers struct {
	ID              uint   `json:"id"`
	Username        string `json:"username"`
	FullName 		string `json:"fullname"`
	ProfilePicture string `json:"profile_picture"`
}

type FormatLikes struct {
	ID   uint        `json:"id"`
	User FormatUsers `json:"user"`
}

// Private Function
func formatUser(user *models.User) FormatUsers {
	return FormatUsers{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
	}
}

func formatLikes(likes []models.Likes) []FormatLikes {
	formattedLikes := make([]FormatLikes, len(likes))
	for i, like := range likes {
		formattedLikes[i] = FormatLikes{
			ID:   like.ID,
			User: formatUser(like.User),
		}
	}
	return formattedLikes
}

func formatComments(comments []models.Comment) []FormatComments {
	formattedComments := make([]FormatComments, len(comments))
	for i, comment := range comments {
		formattedComments[i] = FormatComments{
			ID:          comment.ID,
			Description: comment.Description,
			CreatedAt:   comment.CreatedAt,
			User:        formatUser(comment.User),
			Likes:       formatLikes(comment.Likes),
		}
	}
	return formattedComments
}

func formatPosts(posts []models.Post) []FormatPosts {
	formattedPosts := make([]FormatPosts, len(posts))
	for i, post := range posts {
		formattedPosts[i] = FormatPosts{
			ID:        post.ID,
			Post:      post.Post,
			CreatedAt: post.CreatedAt,
			User:      formatUser(post.User),
			Likes:     formatLikes(post.Likes),
			Comment:   formatComments(post.Comment),
		}
	}
	return formattedPosts
}

func formatOnePost(post models.Post) FormatPosts {
	formattedPost := FormatPosts{
		ID:        post.ID,
		Post:      post.Post,
		CreatedAt: post.CreatedAt,
		User:      formatUser(post.User),
		Likes:     formatLikes(post.Likes),
		Comment:   formatComments(post.Comment),
	}
	return formattedPost
}

// Function
func CreatePost(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	Post := models.Post{}

	if contentType == appJSON {
		c.ShouldBindJSON(&Post)
	} else {
		c.ShouldBind(&Post)
	}

	Post.UserID = userID

	err := db.Debug().Create(&Post).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":        Post.ID,
			"post":      Post.Post,
			"createdAt": Post.CreatedAt,
			"user_id":   Post.UserID,
		},
	})
}

func GetPosts(c *gin.Context) {
	db := database.GetDB()

	var posts []models.Post

	if err := db.Debug().
		Preload("User").
		Preload("Likes").
		Preload("Likes.User").
		Preload("Comment", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Preload("Likes").Preload("Likes.User")
		}).
		Order("created_at desc").
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	response := formatPosts(posts)

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func GetPostByID(c *gin.Context) {
	db := database.GetDB()

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

	var post models.Post

	if err := db.Debug().
		Preload("User").
		Preload("Likes").
		Preload("Likes.User").
		Preload("Comment", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Preload("Likes").Preload("Likes.User")
		}).
		Where("id=?", postIDUint).
		Take(&post).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	formattedPost := formatOnePost(post)

	c.JSON(http.StatusOK, gin.H{
		"data": formattedPost,
	})
}

func GetPostByUserID(c *gin.Context) {
	db := database.GetDB()

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

	var posts []models.Post

	if err := db.Debug().
		Preload("User").
		Preload("Likes").
		Preload("Likes.User").
		Preload("Comment", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User").Preload("Likes").Preload("Likes.User")
		}).
		Where("user_id=?", userIDUint).
		Order("created_at desc").
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	response := formatPosts(posts)

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func DeletePost(c *gin.Context) {
	db := database.GetDB()

	var post models.Post

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

	if err := db.Debug().Where("id=?", postIDUint).
		Where("user_id=?", userID).
		Take(&post).
		Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	if err := db.Debug().
		Delete(&post).
		Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Post with id %d deleted successfully", postIDUint),
	})
}
