package dto

import (
	"time"

	"user-service/internal/models"
)

type CreateUserProfileReq struct {
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Phone       string         `json:"phone"`
	DateOfBirth time.Time      `json:"date_of_birth"`
	Gender      models.Gender  `json:"gender"`
	Bio         string         `json:"bio"`
	Address     models.Address `json:"address"`
}
