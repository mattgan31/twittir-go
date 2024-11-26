package models

import (
	"twittir-go/helpers"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName       string         `gorm:"not null" json:"full_name" validate:"required"`
	Username       string         `gorm:"not null;unique" json:"username" validate:"required"`
	Email          string         `gorm:"not null;unique" json:"email" validate:"required,email"`
	Password       string         `gorm:"not null" json:"-" validate:"required,min=6"`
	Bio            string         `gorm:"type:text;default:null" json:"bio" validate:"max=255"`
	ProfilePicture string         `gorm:"type:text;default:null" json:"profile_picture"`
	SearchVector   string         `gorm:"type:tsvector;generated" json:"-"`
	Posts          []Post         `gorm:"foreignKey:UserID" json:"posts"`
	Comments       []Comment      `gorm:"foreignKey:UserID" json:"comments"`
	Likes          []Likes        `gorm:"foreignKey:UserID" json:"likes"`
	Followers      []Relationship `gorm:"foreignKey:FollowerID" json:"followers"`
	Following      []Relationship `gorm:"foreignKey:FollowingID" json:"following"`
}

func (u *User) BeforeCreate(g *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = helpers.HashPass(u.Password)
	err = nil
	return
}
