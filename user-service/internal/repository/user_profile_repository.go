package repository

import (
	"user-service/internal/models"

	"gorm.io/gorm"
)

type UserProfileRepository interface {
	Create(profile *models.UserProfile) error
	GetByUserID(userID uint) (*models.UserProfile, error)
}

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) Create(profile *models.UserProfile) error {
	return r.db.Create(profile).Error
}

func (r *userProfileRepository) GetByUserID(userID uint) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := r.db.First(&profile, userID).Error
	return &profile, err
}
