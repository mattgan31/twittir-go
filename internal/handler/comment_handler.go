package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"twittir-go/internal/database"
	"twittir-go/internal/domain"
	"twittir-go/internal/helpers"
	"twittir-go/internal/services"
	"twittir-go/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type CommentHandler struct {
	commentService services.CommentService
	userHandler    *UserHandler
}

func NewCommentHandler(commentService services.CommentService, userHandler *UserHandler) *CommentHandler {
	return &CommentHandler{commentService, userHandler}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {

	contentType := helpers.GetContentType(c)

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	postIDUint := uint(postID)
	Comment := domain.Comment{}

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.CreatedAt = time.Now()
	Comment.UpdatedAt = time.Now()
	Comment.UserID = userID
	Comment.PostID = postIDUint

	comment, err := h.commentService.CreateComment(Comment.Description, Comment.PostID, Comment.UserID)
	// err = db.Debug().Create(&Comment).Error

	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
	db := database.GetDB()

	var comment domain.Comment

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
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
