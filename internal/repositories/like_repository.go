package repositories

import (
	"twittir-go/internal/domain"

	"gorm.io/gorm"
)

type LikeRepository interface {
	LikePost(postID *uint, userID uint) (*domain.Like, error)
	CheckLikePost(postID *uint, userID uint) *domain.Like
	UnlikePost(postID *uint, userID uint) error
	LikeComment(postID *uint, userID uint) (*domain.Like, error)
	CheckLikeComment(postID *uint, userID uint) *domain.Like
	UnlikeComment(postID *uint, userID uint) error
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db}
}

func (r *likeRepository) LikePost(postID *uint, userID uint) (*domain.Like, error) {
	like := &domain.Like{
		PostID: postID,
		UserID: userID,
	}

	err := r.db.Debug().Create(&like).Error
	if err != nil {
		return nil, err
	}

	return like, nil

}

func (r *likeRepository) CheckLikePost(postID *uint, userID uint) *domain.Like {
	var like *domain.Like

	r.db.Debug().Where("post_id = ? AND user_id = ?", *postID, userID).Take(&like)

	if like.ID == 0 {
		return nil
	}

	return like
}

func (r *likeRepository) UnlikePost(postID *uint, userID uint) error {
	return r.db.Debug().Where("post_id = ? AND user_id = ?", postID, userID).Delete(&domain.Like{}).Error

}

func (r *likeRepository) LikeComment(postID *uint, userID uint) (*domain.Like, error) {
	like := &domain.Like{
		CommentID: postID,
		UserID:    userID,
	}

	err := r.db.Debug().Create(&like).Error
	if err != nil {
		return nil, err
	}

	return like, nil

}

func (r *likeRepository) CheckLikeComment(postID *uint, userID uint) *domain.Like {
	var like *domain.Like

	r.db.Debug().Where("comment_id = ? AND user_id = ?", *postID, userID).Take(&like)

	if like.ID == 0 {
		return nil
	}

	return like
}

func (r *likeRepository) UnlikeComment(postID *uint, userID uint) error {
	return r.db.Debug().Where("comment_id = ? AND user_id = ?", postID, userID).Delete(&domain.Like{}).Error

}
