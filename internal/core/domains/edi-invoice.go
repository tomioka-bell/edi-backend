package domains

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

// =============================================== HEADER =======================================================================
type EDIInvoice struct {
	EDIInvoiceID  mssql.UniqueIdentifier `gorm:"column:edi_invoice_id;type:uniqueidentifier;primaryKey"`
	NumberInvoice string                 `gorm:"column:number_invoice;type:nvarchar(50);uniqueIndex"`
	NumberOrder   string                 `gorm:"column:number_order;type:nvarchar(50);index"`
	// ProductCode    *string                `gorm:"column:product_code;type:nvarchar(50)"` ลบ
	VanderCode            string                  `gorm:"column:vendor_code;type:nvarchar(50)"`
	InvoiceType           string                  `gorm:"column:invoice_type;type:nvarchar(50);index"`
	ReadInvoice           bool                    `gorm:"column:read_invoice;default:false"`
	StatusInvoice         string                  `gorm:"column:status_invoice;type:nvarchar(30);index"`
	FileURL               *string                 `gorm:"column:file_url;type:nvarchar(500)"`
	ActiveVersionID       *mssql.UniqueIdentifier `gorm:"column:active_version_id;type:uniqueidentifier"`
	VendorCode            string                  `gorm:"column:vendor_code;type:nvarchar(50);index"`
	CreatedByExternalID   string                  `gorm:"column:created_by_external_id;type:nvarchar(50);index"`
	CreatedBySourceSystem string                  `gorm:"column:created_by_source_system;type:nvarchar(50);index"`
	CreatedByPrincipal    *EDI_Principal          `gorm:"foreignKey:CreatedByExternalID,CreatedBySourceSystem;references:ExternalID,SourceSystem"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetimeoffset(7);autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetimeoffset(7);index"`

	RowVer   []byte              `gorm:"column:row_ver;->;-:migration"`
	Versions []EDIInvoiceVersion `gorm:"foreignKey:EDIInvoiceID"`
}

func (EDIInvoice) TableName() string { return "edi_invoice" }

/* ================================================= VERSION ================================================================ */

type EDIInvoiceVersion struct {
	EDIInvoiceVersionID mssql.UniqueIdentifier `gorm:"column:edi_invoice_version_id;type:uniqueidentifier;primaryKey"`
	EDIInvoiceID        mssql.UniqueIdentifier `gorm:"column:edi_invoice_id;type:uniqueidentifier;index;uniqueIndex:ux_efv_doc_ver"`
	VersionNo           int                    `gorm:"column:version_no;not null;uniqueIndex:ux_efv_doc_ver"`
	PeriodFrom          *time.Time             `gorm:"column:period_from;type:datetimeoffset(7);index"`
	PeriodTo            *time.Time             `gorm:"column:period_to;type:datetimeoffset(7);index"`
	StatusInvoice       string                 `gorm:"column:status_invoice;type:nvarchar(30);index"`
	Note                string                 `gorm:"column:note;type:nvarchar(500)"`
	SourceFileURL       *string                `gorm:"column:source_file_url;type:nvarchar(500)"`

	CreatedByExternalID   string         `gorm:"column:created_by_external_id;type:nvarchar(50);index"`
	CreatedBySourceSystem string         `gorm:"column:created_by_source_system;type:nvarchar(50);index"`
	CreatedByPrincipal    *EDI_Principal `gorm:"foreignKey:CreatedByExternalID,CreatedBySourceSystem;references:ExternalID,SourceSystem"`

	//แก้ไขชั่วคราว
	VendorCode    string `json:"vendor_code"`
	NumberInvoice string `json:"number_invoice"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetimeoffset(7);autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetimeoffset(7);index"`

	RowVer []byte `gorm:"column:row_ver;->;-:migration"`
}

func (EDIInvoiceVersion) TableName() string { return "edi_invoice_version" }

// ============================================= STATUS LOG ====================================================================

type EDIInvoiceVersionStatusLog struct {
	EDIInvoiceVersionStatusLogID mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;primaryKey;default:newid()" json:"edi_invoice_version_status_log_id"`
	EDIInvoiceID                 mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;index" json:"edi_invoice_id"`

	OldStatus string `gorm:"type:nvarchar(30)" json:"old_status"`
	NewStatus string `gorm:"type:nvarchar(30);index" json:"new_status"`
	Note      string `gorm:"type:nvarchar(1000)" json:"note"`

	FileURL *string `gorm:"type:nvarchar(500)" json:"file_url"`

	ChangedByExternalID   string         `gorm:"column:changed_by_external_id;type:nvarchar(50);index"`
	ChangedBySourceSystem string         `gorm:"column:changed_by_source_system;type:nvarchar(50);index"`
	ChangedByPrincipal    *EDI_Principal `gorm:"foreignKey:ChangedByExternalID,ChangedBySourceSystem;references:ExternalID,SourceSystem"`

	CreatedAt time.Time `gorm:"type:datetimeoffset(7);autoCreateTime" json:"created_at"`
}

func (EDIInvoiceVersionStatusLog) TableName() string { return "edi_invoice_version_status_log" }

// ============================================= ACTIVE ====================================================================

type EDIInvoiceWithActive struct {
	// header
	EDIInvoiceID    mssql.UniqueIdentifier  `gorm:"column:edi_invoice_id"`
	NumberInvoice   string                  `gorm:"column:number_invoice"`
	VendorCode      string                  `gorm:"column:vendor_code"`
	NumberOrder     string                  `gorm:"column:number_order"`
	ReadInvoice     bool                    `gorm:"column:read_invoice"`
	StatusInvoice   string                  `gorm:"column:status_invoice"`
	FileURL         *string                 `gorm:"column:file_url"`
	CreatedAt       time.Time               `gorm:"column:created_at"`
	UpdatedAt       time.Time               `gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt          `gorm:"column:deleted_at"`
	RowVer          []byte                  `gorm:"column:row_ver"`
	ActiveVersionID *mssql.UniqueIdentifier `gorm:"column:active_version_id"`

	// active version
	AV_ID            mssql.UniqueIdentifier `gorm:"column:av_id"`
	AV_VersionNo     int                    `gorm:"column:av_version_no"`
	AV_PeriodFrom    *time.Time             `gorm:"column:av_period_from"`
	AV_PeriodTo      *time.Time             `gorm:"column:av_period_to"`
	AV_Status        string                 `gorm:"column:av_status"`
	AV_Read          bool                   `gorm:"column:av_read"`
	AV_Note          string                 `gorm:"column:av_note"`
	AV_SourceFileURL *string                `gorm:"column:av_source_file_url"`
	AV_CreatedAt     time.Time              `gorm:"column:av_created_at"`
	AV_DeletedAt     gorm.DeletedAt         `gorm:"column:av_deleted_at"`
	AV_RowVer        []byte                 `gorm:"column:av_row_ver"`

	// latest status log of active version
	LastStatusLogID *mssql.UniqueIdentifier `gorm:"column:last_status_log_id"`
	LastOldStatus   *string                 `gorm:"column:last_old_status"`
	LastNewStatus   *string                 `gorm:"column:last_new_status"`
	LastStatusNote  *string                 `gorm:"column:last_status_note"`
	LastFileURL     *string                 `gorm:"column:last_file_url"`
	LastStatusAt    *time.Time              `gorm:"column:last_status_at"`
}

// ====================================================================================================================

type InvoiceStatusSummary struct {
	VendorCode    string `gorm:"column:vendor_code"`
	NewCount      int64  `gorm:"column:new_count"`
	ConfirmCount  int64  `gorm:"column:confirm_count"`
	RejectCount   int64  `gorm:"column:reject_count"`
	ApprovedCount int64  `gorm:"column:approved_count"`
	TotalCount    int64  `gorm:"column:total_count"`
}

type InvoiceBasicInfo struct {
	NumberInvoice string `gorm:"column:number_invoice"`
	VendorCode    string `gorm:"column:vendor_code"`
}
