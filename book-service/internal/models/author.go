package models

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null;size:255" json:"name"`
	Bio       string         `gorm:"type:text" json:"bio"`
	BirthDate *time.Time     `json:"birth_date,omitempty"`
	Country   string         `gorm:"size:100" json:"country"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// One-to-many relationship with Books
	Books []Book `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE" json:"books,omitempty"`
}

func (Author) TableName() string {
	return "authors"
}
