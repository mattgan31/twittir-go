package models

import (
	"time"

	"gorm.io/gorm"
)

type GormModel struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Created_At time.Time      `json:"createdAt,omitempty"`
	Updated_At time.Time      `json:"updatedAt,omitempty"`
	Deleted_At gorm.DeletedAt `gorm:"index"`
}
