package models

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

type EDIReadStatResp struct {
	EDIReadID  mssql.UniqueIdentifier `json:"edi_read_id"`
	Number     string                 `json:"number"`
	Type       string                 `json:"type"`
	VendorCode string                 `json:"vendor_code"`
	Read       bool                   `json:"read"`
	ReadAt     time.Time              `json:"read_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type EDIReadStatReq struct {
	EDIReadID  mssql.UniqueIdentifier `json:"edi_read_id"`
	Number     string                 `json:"number"`
	Type       string                 `json:"type"`
	VendorCode string                 `json:"vendor_code"`
	Read       bool                   `json:"read"`
	ReadAt     time.Time              `json:"read_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
