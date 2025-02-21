package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"twittir-go/internal/domain"
	"twittir-go/internal/dto"
	"twittir-go/internal/helpers"
	"twittir-go/internal/services"
	"twittir-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService services.PostService
	userHandler *UserHandler
}

func NewPostHandler(postService services.PostService, userHandler *UserHandler) *PostHandler {
	return &PostHandler{postService, userHandler}
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

// GetPostByID godoc
// @Summary Show post by ID
// @Description Retrieve the post details by the specified post ID
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.SuccessResponse{data=dto.ProfileResponse} "Post successfully retrieved"
// @Failure 400 {object} dto.ErrorResponse "Invalid request or missing parameters"
// @Failure 404 {object} dto.ErrorResponse "Post not found"
// @Router /api/posts/{id} [get]
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

// GetPostByUserID godoc
// @Summary Show post by ID
// @Description Retrieve the post details by the specified post ID
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.SuccessResponse{data=dto.ProfileResponse} "Post successfully retrieved"
// @Failure 400 {object} dto.ErrorResponse "Invalid request or missing parameters"
// @Failure 404 {object} dto.ErrorResponse "Post not found"
// @Router /api/posts/{id} [get]
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

func (h *PostHandler) GetPostByFollowingUser(c *gin.Context) {

	userID, err := h.userHandler.extractUserID(c)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid token")
		return
	}

	userIDUint := uint(userID)

	var posts []dto.FormatPost

	posts, err = h.postService.GetPostByFollowingUser(userIDUint)

	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, posts)
}

func (h *PostHandler) DeletePost(c *gin.Context) {

	userID, err := h.userHandler.extractUserID(c)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid token")
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.postService.GetPostByID(int(postID))
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if post.User.ID != userID {
		utils.RespondWithError(c, http.StatusUnauthorized, "You are not authorized to delete this post")
		return
	}

	deletePost := h.postService.DeletePost(int(post.ID))
	if deletePost != nil {
		utils.RespondWithError(c, http.StatusBadRequest, deletePost.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, fmt.Sprintf("Post with id %d deleted successfully", postID))
}
