package models

import (
	"time"

	"gorm.io/gorm"
)

type Credential struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string         `gorm:"not null;size:255" json:"-"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	RefreshTokens []RefreshToken `gorm:"foreignKey:CredentialID;constraint:OnDelete:CASCADE" json:"-"`
}

type RefreshToken struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CredentialID uint           `gorm:"not null;index" json:"credential_id"`
	Token        string         `gorm:"uniqueIndex;not null;size:500" json:"token"`
	ExpiresAt    time.Time      `gorm:"not null" json:"expires_at"`
	IsRevoked    bool           `gorm:"default:false" json:"is_revoked"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Credential Credential `gorm:"foreignKey:CredentialID" json:"-"`
}

func (Credential) TableName() string {
	return "credentials"
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
