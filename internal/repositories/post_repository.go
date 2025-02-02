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
	CreatePost(post string, userID uint) (*domain.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db}
}

func (r *postRepository) FindAllPosts() ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.Debug().Preload("User").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) FindPostByID(id int) (*domain.Post, error) {
	var post domain.Post
	err := r.db.Debug().Preload("User").Where("id = ?", id).Take(&post).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) FindPostByUserID(userID uint) ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.Debug().Preload("User").Where("user_id = ?", userID).Find(&posts).Error
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
