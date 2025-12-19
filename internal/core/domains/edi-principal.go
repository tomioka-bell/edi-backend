package domains

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	gorm "gorm.io/gorm"
)

type EDI_Principal struct {
	EDI_PrincipalID mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;primaryKey;default:(newid())" json:"edi_principal_id"`

	ExternalID   string `gorm:"column:external_id;type:nvarchar(50);index;uniqueIndex:ux_principal_ext_src"`
	SourceSystem string `gorm:"column:source_system;type:nvarchar(50);index;uniqueIndex:ux_principal_ext_src"`

	Email           string `gorm:"column:email;type:nvarchar(100);index"`
	DisplayName     string `gorm:"column:display_name;type:nvarchar(100);index"`
	Profile         string `gorm:"column:profile;type:nvarchar(255)"`
	Group           string `gorm:"column:group;type:nvarchar(50)"`
	CompanyCode     string `gorm:"column:company_code;type:nvarchar(50)"`
	Role            string `gorm:"column:role;type:nvarchar(50)"`
	Status          string `gorm:"column:status;type:nvarchar(50);default:active"`
	Username        string `gorm:"column:username;type:nvarchar(50);uniqueIndex"`
	LoginWithoutOTP bool   `gorm:"column:login_without_otp;type:bit;default:0"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetimeoffset(7);autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetimeoffset(7);index"`
}
