package models

import (
	// "twittir-go/helpers"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Relationship struct {
	gorm.Model
	FollowerID  uint `gorm:"not null" json:"follower_id" valid:"required"`
	FollowingID uint `gorm:"not null" json:"following_id" valid:"required"`
}

func (r *Relationship) BeforeCreate(g *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(r)

	if errCreate != nil {
		err = errCreate
		return
	}

	return
}
