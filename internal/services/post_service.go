package services

import (
	"errors"
	"twittir-go/internal/domain"
	"twittir-go/internal/dto"
	"twittir-go/internal/repositories"
)

type PostService interface {
	GetAllPosts() ([]dto.FormatPost, error)
	GetPostByID(id int) (*dto.FormatPost, error)
	GetPostByUserID(userID uint) ([]dto.FormatPost, error)
	GetPostByFollowingUser(userID uint) ([]dto.FormatPost, error)

	CreatePost(post string, userID uint) (*domain.Post, error)

	DeletePost(id int) error
}

type postService struct {
	postRepo repositories.PostRepository
}

func NewPostService(postRepo repositories.PostRepository) PostService {
	return &postService{postRepo}
}

func (s *postService) GetAllPosts() ([]dto.FormatPost, error) {
	posts, err := s.postRepo.FindAllPosts()
	if err != nil {
		return nil, err
	}

	return dto.NewFormatPosts(posts), nil
}

func (s *postService) GetPostByID(id int) (*dto.FormatPost, error) {
	post, err := s.postRepo.FindPostByID(id)
	if err != nil {
		return nil, err
	}

	if post == nil {
		return nil, errors.New("post not found")
	}

	formattedPost := dto.NewFormatPost(*post)
	return &formattedPost, nil
}

func (s *postService) GetPostByUserID(userID uint) ([]dto.FormatPost, error) {
	posts, err := s.postRepo.FindPostByUserID(userID)
	if err != nil {
		return nil, err
	}

	return dto.NewFormatPosts(posts), nil
}

func (s *postService) GetPostByFollowingUser(userID uint) ([]dto.FormatPost, error) {
	posts, err := s.postRepo.FindPostByFollowingUser(userID)
	if err != nil {
		return nil, err
	}

	return dto.NewFormatPosts(posts), nil
}

// CREATE POST
func (s *postService) CreatePost(post string, userID uint) (*domain.Post, error) {
	// Create a new post

	createdPost, err := s.postRepo.CreatePost(post, userID)
	if err != nil {
		return nil, err
	}

	return createdPost, nil
}

func (s *postService) DeletePost(id int) error {
	err := s.postRepo.DeletePost(id)
	if err != nil {
		return err
	}
	return nil
}
