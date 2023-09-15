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

type FormatPosts struct {
	ID         uint             `json:"id"`
	Post       string           `json:"post"`
	Created_At time.Time        `json:"createdAt"`
	User       FormatUsers      `json:"user"`
	Likes      []FormatLikes    `json:"likes"`
	Comment    []FormatComments `json:"comments"`
}

type FormatComments struct {
	ID          uint          `json:"id"`
	Description string        `json:"description"`
	Created_At  time.Time     `json:"createdAt"`
	User        FormatUsers   `json:"user"`
	Likes       []FormatLikes `json:"likes"`
}

type FormatUsers struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type FormatLikes struct {
	ID   uint        `json:"id"`
	User FormatUsers `json:"user"`
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
			"id":        Post.ID,
			"post":      Post.Post,
			"createdAt": Post.Created_At,
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

func GetPostByID(c *gin.Context) {
	db := database.GetDB()

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PostID"})
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
		Find(&post).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{
		"posts": formattedPost,
	})
}

func GetPostByUserID(c *gin.Context) {
	db := database.GetDB()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
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

func DeletePost(c *gin.Context) {
	db := database.GetDB()

	var post models.Post

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PostID"})
		return
	}

	postIDUint := uint(postID)

	if err := db.Debug().Where("id=?", postIDUint).
		Where("user_id=?", userID).
		First(&post).
		Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := db.Debug().
		Delete(&post).
		Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Post with id %d deleted successfully", postIDUint),
	})
}
