package dto

import (
	"time"
)

type CreateAuthorReq struct {
	Name      string     `json:"name" binding:"required"`
	Bio       string     `json:"bio"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Country   string     `json:"country"`
}

type UpdateAuthorReq struct {
	Name      *string    `json:"name,omitempty"`
	Bio       *string    `json:"bio,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Country   *string    `json:"country,omitempty"`
}

type CreateBookReq struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	PublishYear int     `json:"publish_year" binding:"required"`
	ISBN        string  `json:"isbn"`
	Genre       string  `json:"genre"`
	Pages       int     `json:"pages"`
	Price       float64 `json:"price"`
	AuthorID    uint    `json:"author_id" binding:"required"`
}

type UpdateBookReq struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	PublishYear *int     `json:"publish_year,omitempty"`
	ISBN        *string  `json:"isbn,omitempty"`
	Genre       *string  `json:"genre,omitempty"`
	Pages       *int     `json:"pages,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	AuthorID    *uint    `json:"author_id,omitempty"`
}

type SearchBooksReq struct {
	AuthorName  string `json:"author_name,omitempty" form:"author_name"`
	BookTitle   string `json:"book_title,omitempty" form:"book_title"`
	PublishYear int    `json:"publish_year,omitempty" form:"publish_year"`
	Genre       string `json:"genre,omitempty" form:"genre"`
	Page        int    `json:"page,omitempty" form:"page"`
	Limit       int    `json:"limit,omitempty" form:"limit"`
}

type SearchBooksRes struct {
	Books      []BookWithAuthor `json:"books"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

type BookWithAuthor struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublishYear int       `json:"publish_year"`
	ISBN        string    `json:"isbn"`
	Genre       string    `json:"genre"`
	Pages       int       `json:"pages"`
	Price       float64   `json:"price"`
	AuthorID    uint      `json:"author_id"`
	AuthorName  string    `json:"author_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
