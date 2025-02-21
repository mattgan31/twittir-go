package services

import (
	"fmt"
	"twittir-go/internal/repositories"
)

type LikeService interface {
	LikePost(postID *uint, userID uint) (string, error)
	LikeComment(postID *uint, userID uint) (string, error)
}

type likeService struct {
	likeRepo repositories.LikeRepository
}

func NewLikeService(likeRepo repositories.LikeRepository) LikeService {
	return &likeService{likeRepo}
}

func (s *likeService) LikePost(postID *uint, userID uint) (string, error) {

	checkLiked := s.likeRepo.CheckLikePost(postID, userID)

	if checkLiked != nil {
		err := s.likeRepo.UnlikePost(postID, userID)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Post with id %d unliked successfully", *postID), nil
	}

	_, err := s.likeRepo.LikePost(postID, userID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Post with id %d liked successfully", *postID), nil
}

func (s *likeService) LikeComment(postID *uint, userID uint) (string, error) {

	checkLiked := s.likeRepo.CheckLikeComment(postID, userID)

	if checkLiked != nil {
		err := s.likeRepo.UnlikeComment(postID, userID)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Comment with id %d unliked successfully", *postID), nil
	}

	_, err := s.likeRepo.LikeComment(postID, userID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Comment with id %d liked successfully", *postID), nil
}
