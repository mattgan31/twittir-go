package dto

import "twittir-go/internal/domain"

type FormatLike struct {
	ID   uint        `json:"id"`
	User FormatUsers `json:"user"`
}

func NewFormatLikes(likes []domain.Like) []FormatLike {
	formattedLikes := make([]FormatLike, len(likes))
	for i, like := range likes {
		var formattedUser FormatUsers

		if like.User != nil { // Nil check
			formattedUser = FormatUsers{
				ID:       like.User.ID,
				Username: like.User.Username,
				FullName: like.User.FullName,
			}
		}

		formattedLikes[i] = FormatLike{
			ID:   like.ID,
			User: formattedUser,
		}
	}
	return formattedLikes
}
