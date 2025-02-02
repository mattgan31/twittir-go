package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"twittir-go/internal/database"
	"twittir-go/internal/domain"
	"twittir-go/internal/dto"
	"twittir-go/internal/helpers"
	"twittir-go/internal/services"
	"twittir-go/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type PostHandler struct {
	postService services.PostService
	userHandler *UserHandler
}

func NewPostHandler(postService services.PostService, userHandler *UserHandler) *PostHandler {
	return &PostHandler{postService, userHandler}
}

// Private Function
func formatUser(user *domain.User) dto.FormatUsers {
	return dto.FormatUsers{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
	}
}

func formatLikes(likes []domain.Likes) []dto.FormatLike {
	formattedLikes := make([]dto.FormatLike, len(likes))
	for i, like := range likes {
		formattedLikes[i] = dto.FormatLike{
			ID:   like.ID,
			User: formatUser(like.User),
		}
	}
	return formattedLikes
}

func formatComments(comments []domain.Comment) []dto.FormatComment {
	formattedComments := make([]dto.FormatComment, len(comments))
	for i, comment := range comments {
		formattedComments[i] = dto.FormatComment{
			ID:          comment.ID,
			Description: comment.Description,
			CreatedAt:   comment.CreatedAt,
			User:        formatUser(comment.User),
			Likes:       formatLikes(comment.Likes),
		}
	}
	return formattedComments
}

func formatPosts(posts []domain.Post) []dto.FormatPost {
	formattedPosts := make([]dto.FormatPost, len(posts))
	for i, post := range posts {
		formattedPosts[i] = dto.FormatPost{
			ID:        post.ID,
			Post:      post.Post,
			CreatedAt: post.CreatedAt,
			User:      formatUser(post.User),
			Like:      formatLikes(post.Likes),
			Comment:   formatComments(post.Comment),
		}
	}
	return formattedPosts
}

func formatOnePost(post domain.Post) dto.FormatPost {
	formattedPost := dto.FormatPost{
		ID:        post.ID,
		Post:      post.Post,
		CreatedAt: post.CreatedAt,
		User:      formatUser(post.User),
		Like:      formatLikes(post.Likes),
		Comment:   formatComments(post.Comment),
	}
	return formattedPost
}

// Function
func (h *PostHandler) CreatePost(c *gin.Context) {

	contentType := helpers.GetContentType(c)

	userID, err := h.userHandler.extractUserID(c)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid token")
		return
	}

	Post := domain.Post{}

	if contentType == appJSON {
		c.ShouldBindJSON(&Post)
	} else {
		c.ShouldBind(&Post)
	}

	Post.UserID = userID

	// err := db.Debug().Create(&Post).Error
	createdPost, err := h.postService.CreatePost(Post.Post, Post.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	formattedPost := dto.FormatCreatedPost{
		ID:        createdPost.ID,
		Post:      createdPost.Post,
		CreatedAt: createdPost.CreatedAt.String(),
	}

	utils.RespondWithSuccess(c, http.StatusOK, formattedPost)
}

func (h *PostHandler) GetPosts(c *gin.Context) {

	posts, err := h.postService.GetAllPosts()
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, posts)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "BAD_REQUEST",
			"error":  "Invalid PostID",
		})
		return
	}

	postIDInt := int(postID)

	post, err := h.postService.GetPostByID(postIDInt)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, post)
}

func (h *PostHandler) GetPostByUserID(c *gin.Context) {

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	userIDUint := uint(userID)

	var posts []dto.FormatPost

	posts, err = h.postService.GetPostByUserID(userIDUint)

	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, posts)
}

func DeletePost(c *gin.Context) {
	db := database.GetDB()

	var post domain.Post

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	postIDUint := uint(postID)

	if err := db.Debug().Where("id=?", postIDUint).
		Where("user_id=?", userID).
		Take(&post).
		Error; err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.Debug().
		Delete(&post).
		Error; err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Post with id %d deleted successfully", postIDUint),
	})
}
