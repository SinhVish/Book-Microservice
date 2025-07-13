package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null;size:255" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	PublishYear int            `gorm:"not null" json:"publish_year"`
	ISBN        string         `gorm:"uniqueIndex;size:17" json:"isbn"`
	Genre       string         `gorm:"size:100" json:"genre"`
	Pages       int            `json:"pages"`
	Price       float64        `gorm:"type:decimal(10,2)" json:"price"`
	AuthorID    uint           `gorm:"not null" json:"author_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Many-to-one relationship with Author
	Author Author `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

func (Book) TableName() string {
	return "books"
}
