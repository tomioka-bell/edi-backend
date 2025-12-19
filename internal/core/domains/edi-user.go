package domains

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

type User struct {
	UserID          mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;primaryKey" json:"user_id"`
	Firstname       string                 `json:"firstname"`
	Lastname        string                 `json:"lastname" `
	Username        string                 `json:"username"`
	Email           string                 `json:"email"`
	Password        string                 `json:"password"`
	Status          string                 `json:"status"`
	Profile         string                 `json:"profile"`
	Role            string                 `json:"role"`
	TempOTP         string                 `gorm:"-"`
	Group           string                 `json:"group"`
	CompanyCode     string                 `gorm:"column:company_code;type:nvarchar(50)"`
	LoginWithoutOTP bool                   `json:"login_without_otp"`
	CreatedAt       *time.Time             `json:"created_at"`
	UpdatedAt       *time.Time             `json:"updated_at"`
	DeletedAt       gorm.DeletedAt         `gorm:"index"`
}
