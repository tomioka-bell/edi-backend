package domains

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

type PasswordReset struct {
	ID        mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	UserID    mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;index" json:"user_id"`
	TokenHash string                 `gorm:"size:64;uniqueIndex" json:"token_hash"`
	ExpiresAt time.Time              `gorm:"index" json:"expires_at"`
	UsedAt    *time.Time             `gorm:"index" json:"used_at"`
	CreatedAt time.Time              `json:"created_at"`
	DeletedAt gorm.DeletedAt         `gorm:"index" json:"deleted_at"`
}
