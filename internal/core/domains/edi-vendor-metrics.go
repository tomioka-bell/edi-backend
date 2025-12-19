package domains

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

type EDIVendorMetrics struct {
	VendorMetricsID mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;primaryKey;default:(newid())" json:"vendor_metrics_id"`
	Initials        string                 `gorm:"column:initials;type:nvarchar(10)" json:"initials"`
	CompanyName     string                 `gorm:"column:company_name;type:nvarchar(255);index" json:"company_name"`

	ReminderDays int `gorm:"column:reminder_days;type:int;default:1" json:"reminder_days"`

	Active bool `gorm:"column:active;default:true" json:"active"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetimeoffset(7);autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetimeoffset(7);index"`
}

func (EDIVendorMetrics) TableName() string {
	return "edi_vendor_metrics"
}

// สูตรคำนวน OnTimePercent
// OnTimePercent = (จำนวนอ่านทันเวลา / จำนวนเอกสารทั้งหมด) × 100
// ตัวอย่าง:

// มี 10 เอกสาร
// ExpectedResponseTimeHr = 24 ชั่วโมง
// มี 7 ฉบับอ่านภายใน 24 ชั่วโมง
// มี 3 ฉบับอ่านเกิน 24 ชั่วโมง
// OnTimePercent = (7 / 10) × 100 = 70%
