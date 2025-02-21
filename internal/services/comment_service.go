package services

import (
	"twittir-go/internal/dto"
	"twittir-go/internal/repositories"
)

type CommentService interface {
	CreateComment(comment string, postID uint, userID uint) (dto.FormatComment, error)
}

type commentService struct {
	commentRepo repositories.CommentRepository
}

func NewCommentService(commentRepo repositories.CommentRepository) CommentService {
	return &commentService{commentRepo}
}

func (s *commentService) CreateComment(comment string, postID uint, userID uint) (dto.FormatComment, error) {
	newComment, err := s.commentRepo.CreateComment(comment, postID, userID)

	if err != nil {
		return dto.FormatComment{}, err
	}

	return dto.FormatComment{
		ID:          newComment.ID,
		Description: newComment.Description,
		CreatedAt:   newComment.CreatedAt,
		User: dto.FormatUsers{
			ID:             newComment.User.ID,
			Username:       newComment.User.Username,
			ProfilePicture: newComment.User.ProfilePicture,
			FullName:       newComment.User.FullName,
		},
		Likes: dto.NewFormatLikes(newComment.Likes),
	}, nil
}
