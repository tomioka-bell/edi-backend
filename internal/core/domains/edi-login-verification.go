package domains

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
)

type LoginVerification struct {
	ID           mssql.UniqueIdentifier `gorm:"column:id;type:uniqueidentifier;primaryKey"`
	UserID       mssql.UniqueIdentifier `gorm:"column:user_id;type:uniqueidentifier"`
	CodeHash     string                 `gorm:"column:code_hash"`
	ExpiresAt    time.Time              `gorm:"column:expires_at"`
	CreatedAt    time.Time              `gorm:"column:created_at"`
	ConsumedAt   *time.Time             `gorm:"column:consumed_at"`
	AttemptCount int                    `gorm:"column:attempt_count"`
}

func (LoginVerification) TableName() string { return "login_verifications" }
