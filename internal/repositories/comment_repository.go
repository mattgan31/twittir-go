package repositories

import (
	"twittir-go/internal/domain"

	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(comment string, postID uint, userID uint) (*domain.Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db}
}

func (r *commentRepository) CreateComment(comment string, postID uint, userID uint) (*domain.Comment, error) {
	newComment := domain.Comment{
		Description: comment,
		PostID:      postID,
		UserID:      userID,
	}

	if err := r.db.Debug().Create(&newComment).Error; err != nil {
		return nil, err
	}

	if err := r.db.Debug().Preload("User").Preload("Likes").First(&newComment, newComment.ID).Error; err != nil {
		return nil, err
	}

	// Ensure Likes is not nil
	if newComment.Likes == nil {
		newComment.Likes = []domain.Like{}
	} else {
		newComment.Likes = append(newComment.Likes, newComment.Likes...)
	}

	return &newComment, nil
}
