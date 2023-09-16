package models

import (
	"twittir-go/helpers"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Full_Name       string         `gorm:"not null" json:"full_name" valid:"required~Full name is required"`
	Username        string         `gorm:"not null; unique" json:"username" valid:"required~Username is required"`
	Email           string         `gorm:"not null; unique" json:"email" valid:"required~Email is required, email~Invalid email format"`
	Password        string         `gorm:"not null" json:"password" valid:"required~Password is required, stringlength(6|255)~Password has to have a minimum length of 6 characters"`
	Bio             string         `gorm:"not null" json:"bio" valid:"stringlength(0|255)"`
	Profile_Picture string         `gorm:"not null"`
	Post            []Post         `gorm:"foreignKey:UserID"`
	Comment         []Comment      `gorm:"foreignKey:UserID"`
	Likes           []Likes        `gorm:"foreignKey:UserID"`
	Followers       []Relationship `gorm:"foreignKey:FollowerID"`
	Following       []Relationship `gorm:"foreignKey:FollowingID"`
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
