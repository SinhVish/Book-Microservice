package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Status    string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Profile UserProfile `gorm:"foreignKey:ID;references:ID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
}

func (User) TableName() string {
	return "users"
}
