package dto

import (
	"time"
	"twittir-go/internal/domain"
)

// Type for Formatting
type FormatPost struct {
	ID           uint            `json:"id"`
	Post         string          `json:"post"`
	LikeCount    int             `json:"like_count"`
	CommentCount int             `json:"comment_count"`
	CreatedAt    time.Time       `json:"created_at"`
	User         FormatUsers     `json:"user"`
	Like         []FormatLike    `json:"likes"`
	Comment      []FormatComment `json:"comments"`
}

type FormatComment struct {
	ID          uint         `json:"id"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"createdAt"`
	User        FormatUsers  `json:"user"`
	Likes       []FormatLike `json:"likes"`
}

type FormatLike struct {
	ID   uint        `json:"id"`
	User FormatUsers `json:"user"`
}

type FormatCreatedPost struct {
	ID        uint   `json:"id"`
	Post      string `json:"post"`
	CreatedAt string `json:"created_at"`
}

func NewFormatPost(post domain.Post) FormatPost {

	return FormatPost{
		ID:   post.ID,
		Post: post.Post,
		User: FormatUsers{
			ID:       post.User.ID,
			Username: post.User.Username,
			FullName: post.User.FullName,
		},
		CreatedAt:    post.CreatedAt,
		LikeCount:    len(post.Likes),
		CommentCount: len(post.Comment),
	}
}

// Convert `[]domain.Post` → `[]dto.FormatPost`
func NewFormatPosts(posts []domain.Post) []FormatPost {
	formattedPosts := make([]FormatPost, len(posts))
	for i, post := range posts {
		formattedPosts[i] = NewFormatPost(post)
	}
	return formattedPosts
}
