package domain

import (
	// "gorm.io/gorm"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	UserID    uint  `gorm:"not null" json:"user_id" valid:"required"`
	PostID    *uint `gorm:"null" json:"post_id"`
	CommentID *uint `gorm:"null" json:"comment_id"`
	User      *User
}

func (l *Like) BeforeCreate(g *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(l)

	if errCreate != nil {
		err = errCreate
		return
	}

	return
}
