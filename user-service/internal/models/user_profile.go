package models

import (
	"time"

	"gorm.io/gorm"
)

type UserProfile struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	FirstName   *string    `gorm:"size:100" json:"first_name,omitempty"`
	LastName    *string    `gorm:"size:100" json:"last_name,omitempty"`
	Phone       *string    `gorm:"size:20" json:"phone,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Gender      *Gender    `gorm:"type:varchar(15);default:'not_specified'" json:"gender,omitempty"`
	Bio         *string    `gorm:"type:text" json:"bio,omitempty"`

	Address *Address `gorm:"embedded;embeddedPrefix:address_" json:"address,omitempty"`

	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Address struct {
	Street  *string `gorm:"size:255" json:"street,omitempty"`
	City    *string `gorm:"size:100" json:"city,omitempty"`
	State   *string `gorm:"size:100" json:"state,omitempty"`
	ZipCode *string `gorm:"size:20" json:"zip_code,omitempty"`
	Country *string `gorm:"size:100" json:"country,omitempty"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
