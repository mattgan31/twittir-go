package domain

import (
	// "twittir-go/helpers"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID      uint    `gorm:"not null" json:"user_id" valid:"required"`
	PostID      uint    `gorm:"not null" json:"post_id" valid:"required"`
	Description string  `gorm:"not null" json:"description" valid:"required~Description Comment is required"`
	Likes       []Likes `gorm:"foreignKey:CommentID"`
	User        *User
}

func (c *Comment) BeforeCreate(g *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(c)

	if errCreate != nil {
		err = errCreate
		return
	}

	return
}
