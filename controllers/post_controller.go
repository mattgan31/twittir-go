package controllers

import (
	"net/http"
	"time"
	"twittir-go/database"
	"twittir-go/helpers"
	"twittir-go/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type FormatPosts struct {
	ID         uint      `json:"id"`
	Post       string    `json:"post"`
	Created_At time.Time `json:"created_at"`
	User       FormatUsers
	Likes      []FormatLikes
	Comment    []FormatComments
}

type FormatComments struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	Created_At  time.Time `json:"created_at"`
	User        FormatUsers
	Likes       []FormatLikes
}

type FormatUsers struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type FormatLikes struct {
	ID   uint `json:"id"`
	User FormatUsers
}

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

	Post.Created_At = time.Now()
	Post.Updated_At = time.Now()
	Post.UserID = userID

	err := db.Debug().Create(&Post).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": gin.H{
			"id":         Post.ID,
			"post":       Post.Post,
			"created_at": Post.Created_At,
			"user_id":    Post.UserID,
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
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	response := make([]FormatPosts, len(posts))

	for i, post := range posts {
		formattedPost := FormatPosts{
			ID:         post.ID,
			Post:       post.Post,
			Created_At: post.Created_At,
		}
		formattedPost.User = FormatUsers{
			ID:       post.User.ID,
			Username: post.User.Username,
		}

		// Format likes for the post
		likes := make([]FormatLikes, len(post.Likes))
		for k, like := range post.Likes {
			likesResponse := FormatLikes{
				ID: like.ID,
			}

			likesResponse.User = FormatUsers{
				ID:       like.User.ID,
				Username: like.User.Username,
			}

			likes[k] = likesResponse
		}

		formattedPost.Likes = likes

		// Format comments for the post
		comments := make([]FormatComments, len(post.Comment))
		for j, comment := range post.Comment {
			commentResponse := FormatComments{
				ID:          comment.ID,
				Description: comment.Description,
				Created_At:  comment.Created_At,
			}

			commentResponse.User = FormatUsers{
				ID:       comment.User.ID,
				Username: comment.User.Username,
			}

			// Format likes for the comment
			commentLikes := make([]FormatLikes, len(comment.Likes))
			for k, like := range comment.Likes {
				commentlike := FormatLikes{
					ID: like.ID,
				}

				commentlike.User = FormatUsers{
					ID:       like.User.ID,
					Username: like.User.Username,
				}

				commentLikes[k] = commentlike
			}

			commentResponse.Likes = commentLikes

			comments[j] = commentResponse
		}
		formattedPost.Comment = comments

		response[i] = formattedPost
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": response,
	})
}
