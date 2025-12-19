package domains

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

type EDIReadStat struct {
	EDIReadID  mssql.UniqueIdentifier `gorm:"column:edi_read_id;type:uniqueidentifier;primaryKey"`
	Number     string                 `gorm:"column:number;type:nvarchar(50);index:ux_read_doc"`
	Type       string                 `gorm:"column:type;type:nvarchar(50);index:ux_read_doc"`
	VendorCode string                 `gorm:"column:vendor_code;type:nvarchar(50);index:ux_read_doc"`
	Read       bool                   `gorm:"column:read"`
	ReadAt     time.Time              `gorm:"column:read_at;type:datetimeoffset(7);autoCreateTime"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetimeoffset(7);autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetimeoffset(7);index"`
}

func (EDIReadStat) TableName() string { return "edi_read_stat" }
