package handler

import (
	"net/http"
	"strconv"
	"twittir-go/internal/services"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeService services.LikeService
	userHandler *UserHandler
}

func NewLikeHandler(postService services.LikeService, userHandler *UserHandler) *LikeHandler {
	return &LikeHandler{postService, userHandler}
}

func (h *LikeHandler) ToggleLikePost(c *gin.Context) {
	userID, err := h.userHandler.extractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PostID"})
		return
	}

	postIDUint := uint(postID)

	likePost, err := h.likeService.LikePost(&postIDUint, userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": likePost,
	})
}

func (h *LikeHandler) ToggleLikeComment(c *gin.Context) {
	userID, err := h.userHandler.extractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PostID"})
		return
	}

	commentIDUint := uint(commentID)

	likePost, err := h.likeService.LikeComment(&commentIDUint, userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": likePost,
	})
}
