package domains

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

// =============================================== HEADER =======================================================================
type EDIOrder struct {
	EDIOrderID      mssql.UniqueIdentifier  `gorm:"column:edi_order_id;type:uniqueidentifier;primaryKey"`
	NumberOrder     string                  `gorm:"column:number_order;type:nvarchar(50);uniqueIndex"`
	NumberForecast  string                  `gorm:"column:number_forecast;type:nvarchar(50);index"`
	VendorCode      string                  `gorm:"column:vendor_code;type:nvarchar(50);index"`
	ReadOrder       bool                    `gorm:"column:read_order;default:false"`
	StatusOrder     string                  `gorm:"column:status_order;type:nvarchar(30);index"`
	FileURL         *string                 `gorm:"column:file_url;type:nvarchar(500)"`
	ActiveVersionID *mssql.UniqueIdentifier `gorm:"column:active_version_id;type:uniqueidentifier"`

	CreatedByExternalID   string         `gorm:"column:created_by_external_id;type:nvarchar(50);index"`
	CreatedBySourceSystem string         `gorm:"column:created_by_source_system;type:nvarchar(50);index"`
	CreatedByPrincipal    *EDI_Principal `gorm:"foreignKey:CreatedByExternalID,CreatedBySourceSystem;references:ExternalID,SourceSystem"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetimeoffset(7);autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetimeoffset(7);index"`

	RowVer   []byte            `gorm:"column:row_ver;->;-:migration"`
	Versions []EDIOrderVersion `gorm:"foreignKey:EDIOrderID"`
}

func (EDIOrder) TableName() string { return "edi_order" }

/* ================================================= VERSION ================================================================ */

type EDIOrderVersion struct {
	EDIOrderVersionID mssql.UniqueIdentifier `gorm:"column:edi_order_version_id;type:uniqueidentifier;primaryKey"`
	EDIOrderID        mssql.UniqueIdentifier `gorm:"column:edi_order_id;type:uniqueidentifier;index;uniqueIndex:ux_efv_doc_ver"`
	VersionNo         int                    `gorm:"column:version_no;not null;uniqueIndex:ux_efv_doc_ver"`

	PeriodFrom    *time.Time `gorm:"column:period_from;type:datetimeoffset(7);index"`
	PeriodTo      *time.Time `gorm:"column:period_to;type:datetimeoffset(7);index"`
	StatusOrder   string     `gorm:"column:status_order;type:nvarchar(30);index"`
	Note          string     `gorm:"column:note;type:nvarchar(500)"`
	SourceFileURL *string    `gorm:"column:source_file_url;type:nvarchar(500)"`

	CreatedByExternalID   string         `gorm:"column:created_by_external_id;type:nvarchar(50);index"`
	CreatedBySourceSystem string         `gorm:"column:created_by_source_system;type:nvarchar(50);index"`
	CreatedByPrincipal    *EDI_Principal `gorm:"foreignKey:CreatedByExternalID,CreatedBySourceSystem;references:ExternalID,SourceSystem"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetimeoffset(7);autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetimeoffset(7);index"`

	//แก้ไขชั่วคราว
	VendorCode  string `json:"vendor_code"`
	NumberOrder string `json:"number_order"`

	RowVer []byte `gorm:"column:row_ver;->;-:migration"`
}

func (EDIOrderVersion) TableName() string { return "edi_order_version" }

// ============================================= STATUS LOG ====================================================================

type EDIOrderVersionStatusLog struct {
	EDIOrderVersionStatusLogID mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;primaryKey;default:newid()" json:"edi_order_version_status_log_id"`
	EDIOrderID                 mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;index" json:"edi_order_id"`

	OldStatus string `gorm:"type:nvarchar(30)" json:"old_status"`
	NewStatus string `gorm:"type:nvarchar(30);index" json:"new_status"`
	Note      string `gorm:"type:nvarchar(1000)" json:"note"`

	FileURL *string `gorm:"type:nvarchar(500)" json:"file_url"`

	ChangedByExternalID   string         `gorm:"column:changed_by_external_id;type:nvarchar(50);index"`
	ChangedBySourceSystem string         `gorm:"column:changed_by_source_system;type:nvarchar(50);index"`
	ChangedByPrincipal    *EDI_Principal `gorm:"foreignKey:ChangedByExternalID,ChangedBySourceSystem;references:ExternalID,SourceSystem"`

	CreatedAt time.Time `gorm:"type:datetimeoffset(7);autoCreateTime" json:"created_at"`
}

func (EDIOrderVersionStatusLog) TableName() string { return "edi_order_version_status_log" }

// ============================================= ACTIVE ====================================================================

type EDIOrderWithActive struct {
	// header
	EDIOrderID      mssql.UniqueIdentifier  `gorm:"column:edi_order_id"`
	NumberOrder     string                  `gorm:"column:number_order"`
	VendorCode      string                  `gorm:"column:vendor_code"`
	NumberForecast  string                  `gorm:"column:number_forecast"`
	ReadOrder       bool                    `gorm:"column:read_order"`
	StatusOrder     string                  `gorm:"column:status_order"`
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

type OrderStatusSummary struct {
	VendorCode    string `gorm:"column:vendor_code"`
	NewCount      int64  `gorm:"column:new_count"`
	ConfirmCount  int64  `gorm:"column:confirm_count"`
	RejectCount   int64  `gorm:"column:reject_count"`
	ApprovedCount int64  `gorm:"column:approved_count"`
	TotalCount    int64  `gorm:"column:total_count"`
}

type StatusTotalSummary struct {
	NewCount      int64 `json:"new_count"`
	ConfirmCount  int64 `json:"confirm_count"`
	RejectCount   int64 `json:"reject_count"`
	ApprovedCount int64 `json:"approved_count"`
	TotalCount    int64 `json:"total_count"`
}

type EDIOrderHeaderWithPeriod struct {
	EDIOrderID     mssql.UniqueIdentifier `gorm:"column:edi_order_id"`
	NumberOrder    string                 `gorm:"column:number_order"`
	NumberForecast string                 `gorm:"column:number_forecast"`
	VendorCode     string                 `gorm:"column:vendor_code"`
	StatusOrder    string                 `gorm:"column:status_order"`
	CreatedAt      time.Time              `gorm:"column:created_at"`
	PeriodTo       *time.Time             `gorm:"column:period_to"`
}

type OrderBasicInfo struct {
	NumberOrder string `gorm:"column:number_order"`
	VendorCode  string `gorm:"column:vendor_code"`
}
