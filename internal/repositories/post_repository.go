package repositories

import (
	"errors"
	"twittir-go/internal/domain"

	"gorm.io/gorm"
)

type PostRepository interface {
	FindAllPosts() ([]domain.Post, error)
	FindPostByID(id int) (*domain.Post, error)
	FindPostByUserID(userID uint) ([]domain.Post, error)
	FindPostByFollowingUser(userID uint) ([]domain.Post, error)
	CreatePost(post string, userID uint) (*domain.Post, error)
	DeletePost(id int) error
	GetLikeCount(postID uint) (int, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db}
}

func (r *postRepository) FindAllPosts() ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.Debug().
		Preload("User").
		Preload("Comment.User").
		Preload("Comment.Likes.User").
		Preload("Likes.User").
		Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) FindPostByID(id int) (*domain.Post, error) {
	var post domain.Post
	err := r.db.Debug().
		Preload("User").
		Preload("Comment.User").
		Preload("Comment.Likes.User").
		Preload("Likes.User").
		Where("id = ?", id).Take(&post).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	if post.Likes == nil {
		post.Likes = []domain.Like{}
	}

	return &post, nil
}

func (r *postRepository) FindPostByUserID(userID uint) ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.Debug().
		Preload("User").
		Preload("Comment.User").
		Preload("Comment.Likes.User").
		Preload("Likes.User").
		Where("user_id = ?", userID).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) FindPostByFollowingUser(userID uint) ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.Debug().
		Joins("JOIN relationships ON relationships.following_id = posts.user_id").
		Where("relationships.follower_id = ?", userID).
		Preload("User").
		Preload("Comment.User").
		Preload("Comment.Likes.User").
		Preload("Likes.User").
		Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) CreatePost(post string, userID uint) (*domain.Post, error) {
	// Create a new post
	newPost := domain.Post{
		Post:   post,
		UserID: userID,
	}

	err := r.db.Debug().Preload("User").Create(&newPost).Error
	if err != nil {
		return nil, err
	}
	return &newPost, nil
}

func (r *postRepository) DeletePost(id int) error {
	err := r.db.Debug().Where("id = ?", id).Delete(&domain.Post{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *postRepository) GetLikeCount(postID uint) (int, error) {
	var count int64
	err := r.db.Debug().Model(&domain.Like{}).Where("post_id = ?", postID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
