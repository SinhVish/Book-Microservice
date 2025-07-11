package repository

import (
	"auth-service/internal/models"

	"gorm.io/gorm"
)

type CredentialRepository interface {
	Create(credential *models.Credential) error
	GetByEmail(email string) (*models.Credential, error)
	GetByID(id uint) (*models.Credential, error)
	Update(credential *models.Credential) error
	Delete(id uint) error
	CreateRefreshToken(refreshToken *models.RefreshToken) error
	GetRefreshTokenByToken(token string) (*models.RefreshToken, error)
	RevokeRefreshToken(token string) error
	RevokeAllRefreshTokensForCredential(credentialID uint) error
}

type credentialRepository struct {
	db *gorm.DB
}

func NewCredentialRepository(db *gorm.DB) CredentialRepository {
	return &credentialRepository{db: db}
}

func (r *credentialRepository) Create(credential *models.Credential) error {
	return r.db.Create(credential).Error
}

func (r *credentialRepository) GetByEmail(email string) (*models.Credential, error) {
	var credential models.Credential
	err := r.db.Where("email = ?", email).First(&credential).Error
	if err != nil {
		return nil, err
	}
	return &credential, nil
}

func (r *credentialRepository) GetByID(id uint) (*models.Credential, error) {
	var credential models.Credential
	err := r.db.First(&credential, id).Error
	if err != nil {
		return nil, err
	}
	return &credential, nil
}

func (r *credentialRepository) Update(credential *models.Credential) error {
	return r.db.Save(credential).Error
}

func (r *credentialRepository) Delete(id uint) error {
	return r.db.Delete(&models.Credential{}, id).Error
}

func (r *credentialRepository) CreateRefreshToken(refreshToken *models.RefreshToken) error {
	return r.db.Create(refreshToken).Error
}

func (r *credentialRepository) GetRefreshTokenByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("token = ? AND is_revoked = false", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *credentialRepository) RevokeRefreshToken(token string) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true).Error
}

func (r *credentialRepository) RevokeAllRefreshTokensForCredential(credentialID uint) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("credential_id = ?", credentialID).
		Update("is_revoked", true).Error
}
