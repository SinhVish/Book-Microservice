package models

import (
	"time"

	"gorm.io/gorm"
)

type UserProfile struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	FirstName   string     `gorm:"size:100" json:"first_name"`
	LastName    string     `gorm:"size:100" json:"last_name"`
	Phone       string     `gorm:"size:20" json:"phone"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      Gender     `gorm:"type:varchar(15);default:'not_specified'" json:"gender"`
	Bio         string     `gorm:"type:text" json:"bio"`

	Address Address `gorm:"embedded;embeddedPrefix:address_" json:"address"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Address struct {
	Street  string `gorm:"size:255" json:"street"`
	City    string `gorm:"size:100" json:"city"`
	State   string `gorm:"size:100" json:"state"`
	ZipCode string `gorm:"size:20" json:"zip_code"`
	Country string `gorm:"size:100" json:"country"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
