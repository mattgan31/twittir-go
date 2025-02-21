package dto

import (
	"time"
	"twittir-go/internal/domain"
)

type FormatComment struct {
	ID          uint         `json:"id"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"createdAt"`
	User        FormatUsers  `json:"user"`
	Likes       []FormatLike `json:"likes"`
}

func NewFormatComments(comments []domain.Comment) []FormatComment {
	formattedComments := make([]FormatComment, len(comments))
	for i, comment := range comments {
		var formattedUser FormatUsers

		if comment.User != nil { // Nil check
			formattedUser = FormatUsers{
				ID:       comment.User.ID,
				Username: comment.User.Username,
				FullName: comment.User.FullName,
			}
		}

		formattedComments[i] = FormatComment{
			ID:          comment.ID,
			Description: comment.Description,
			CreatedAt:   comment.CreatedAt,
			User:        formattedUser,
			Likes:       NewFormatLikes(comment.Likes),
		}
	}
	return formattedComments
}
