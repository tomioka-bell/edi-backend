package domains

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

// บริษัท (ตัวอย่าง)
type EDICompany struct {
	CompanyID mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;primaryKey;default:(newid())" json:"company_id"`
	Name      string                 `gorm:"column:name;type:nvarchar(200)"`

	// preload ใช้ได้
	NotificationRecipients []EDICompanyNotificationRecipient `gorm:"foreignKey:CompanyID"`
}

func (EDICompany) TableName() string { return "edi_company" }

// ผู้รับอีเมลของบริษัทนั้น ๆ
type EDICompanyNotificationRecipient struct {
	CompanyNotificationRecipientID mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;primaryKey;default:newid()" json:"company_notification_recipient_id"`

	CompanyID mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;index" json:"company_id"`
	// "FORECAST", "ORDER", "INVOICE"
	NotificationType string `gorm:"type:nvarchar(50);index" json:"notification_type"`

	Email string `gorm:"type:nvarchar(255);index" json:"email"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetimeoffset(7);autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetimeoffset(7);index"`
}

func (EDICompanyNotificationRecipient) TableName() string {
	return "edi_company_notification_recipient"
}

// ผู้รับอีเมลของ vendor นั้น ๆ
type EDIVendorNotificationRecipient struct {
	VendorNotificationRecipientID mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;primaryKey;default:newid()" json:"vendor_notification_recipient_id"`
	Company                       string                 `gorm:"type:nvarchar(255);index" json:"company"`
	// "FORECAST", "ORDER", "INVOICE", "VENDOR"
	NotificationType string `gorm:"type:nvarchar(50);index" json:"notification_type"`

	// FK ไปยัง EDI_Principal
	EDI_PrincipalID mssql.UniqueIdentifier `gorm:"column:edi_principal_id;type:uniqueidentifier;index" json:"edi_principal_id"`
	Principal       *EDI_Principal         `gorm:"foreignKey:EDI_PrincipalID;references:EDI_PrincipalID"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetimeoffset(7);autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetimeoffset(7);index"`
}

func (EDIVendorNotificationRecipient) TableName() string {
	return "edi_vendor_notification_recipient"
}
