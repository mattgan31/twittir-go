package domain

import (
	// "twittir-go/helpers"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Post    string    `gorm:"not null" json:"post" valid:"required~Post is required"`
	UserID  uint      `gorm:"not null" json:"user_id" valid:"required"`
	Likes   []Like    `gorm:"foreignKey:PostID"`
	Comment []Comment `gorm:"foreignKey:PostID"`
	User    *User
}

func (p *Post) BeforeCreate(g *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}

	return
}
