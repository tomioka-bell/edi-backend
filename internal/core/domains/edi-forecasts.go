package domains

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

// =============================================== HEADER =======================================================================
type EDI_Forecast struct {
	EDI_ForecastID mssql.UniqueIdentifier `gorm:"column:edi_forecast_id;type:uniqueidentifier;primaryKey"`
	NumberForecast string                 `gorm:"column:number_forecast;type:nvarchar(50);uniqueIndex"`
	VendorCode     string                 `gorm:"column:vendor_code;type:nvarchar(50);index"`
	ReadForecast   bool                   `gorm:"column:read_forecast"`
	ReadAt         *time.Time             `gorm:"column:read_at;type:datetimeoffset(7)"`
	StatusForecast string                 `gorm:"column:status_forecast;type:nvarchar(30);index"`
	FileURL        *string                `gorm:"column:file_url;type:nvarchar(500)"`

	ActiveVersionID *mssql.UniqueIdentifier `gorm:"column:active_version_id;type:uniqueidentifier"`

	CreatedByExternalID   string         `gorm:"column:created_by_external_id;type:nvarchar(50);index"`
	CreatedBySourceSystem string         `gorm:"column:created_by_source_system;type:nvarchar(50);index"`
	CreatedByPrincipal    *EDI_Principal `gorm:"foreignKey:CreatedByExternalID,CreatedBySourceSystem;references:ExternalID,SourceSystem"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetimeoffset(7);autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetimeoffset(7);index"`

	RowVer   []byte                `gorm:"column:row_ver;->;-:migration"`
	Versions []EDI_ForecastVersion `gorm:"foreignKey:EDIForecastID"`
}

func (EDI_Forecast) TableName() string { return "edi_forecast" }

/* ================================================= VERSION ================================================================ */

type EDI_ForecastVersion struct {
	EDIForecastVersionID mssql.UniqueIdentifier `gorm:"column:edi_forecast_version_id;type:uniqueidentifier;primaryKey"`
	EDIForecastID        mssql.UniqueIdentifier `gorm:"column:edi_forecast_id;type:uniqueidentifier;index;uniqueIndex:ux_efv_doc_ver"`
	VersionNo            int                    `gorm:"column:version_no;not null;uniqueIndex:ux_efv_doc_ver"`

	PeriodFrom     *time.Time `gorm:"column:period_from;type:datetimeoffset(7);index"`
	PeriodTo       *time.Time `gorm:"column:period_to;type:datetimeoffset(7);index"`
	StatusForecast string     `gorm:"column:status_forecast;type:nvarchar(30);index"`
	ReadForecast   bool       `gorm:"column:read_forecast"`
	Note           string     `gorm:"column:note;type:nvarchar(500)"`
	SourceFileURL  *string    `gorm:"column:source_file_url;type:nvarchar(500)"`

	CreatedByExternalID   string         `gorm:"column:created_by_external_id;type:nvarchar(50);index"`
	CreatedBySourceSystem string         `gorm:"column:created_by_source_system;type:nvarchar(50);index"`
	CreatedByPrincipal    *EDI_Principal `gorm:"foreignKey:CreatedByExternalID,CreatedBySourceSystem;references:ExternalID,SourceSystem"`

	CreatedAt time.Time      `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetimeoffset(7);autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetimeoffset(7);index"`

	//แก้ไขชั่วคราว
	VendorCode     string `json:"vendor_code"`
	NumberForecast string `json:"number_forecast"`

	RowVer []byte `gorm:"column:row_ver;->;-:migration"`
}

func (EDI_ForecastVersion) TableName() string { return "edi_forecast_version" }

// ============================================= STATUS LOG ====================================================================

type EDI_ForecastVersionStatusLog struct {
	EDIForecastVersionStatusLogID mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;primaryKey;default:newid()" json:"edi_forecast_version_status_log_id"`

	EDIForecastID mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;index" json:"edi_forecast_id"`

	OldStatus string `gorm:"type:nvarchar(30)" json:"old_status"`
	NewStatus string `gorm:"type:nvarchar(30);index" json:"new_status"`
	Note      string `gorm:"type:nvarchar(1000)" json:"note"`

	FileURL *string `gorm:"type:nvarchar(500)" json:"file_url"`

	ChangedByExternalID   string         `gorm:"column:changed_by_external_id;type:nvarchar(50);index"`
	ChangedBySourceSystem string         `gorm:"column:changed_by_source_system;type:nvarchar(50);index"`
	ChangedByPrincipal    *EDI_Principal `gorm:"foreignKey:ChangedByExternalID,ChangedBySourceSystem;references:ExternalID,SourceSystem"`

	CreatedAt time.Time `gorm:"type:datetimeoffset(7);autoCreateTime" json:"created_at"`
}

func (EDI_ForecastVersionStatusLog) TableName() string { return "edi_forecast_version_status_log" }

// ============================================= ACTIVE ====================================================================

type EDIForecastWithActive struct {
	// header
	EDIForecastID   mssql.UniqueIdentifier  `gorm:"column:edi_forecast_id"`
	NumberForecast  string                  `gorm:"column:number_forecast"`
	VendorCode      string                  `gorm:"column:vendor_code"`
	ReadForecast    bool                    `gorm:"column:read_forecast"`
	ActiveVersionID *mssql.UniqueIdentifier `gorm:"column:active_version_id"`
	StatusForecast  string                  `gorm:"column:status_forecast"`
	FileURL         *string                 `gorm:"column:file_url"`
	CreatedAt       time.Time               `gorm:"column:created_at"`
	UpdatedAt       time.Time               `gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt          `gorm:"column:deleted_at"`
	RowVer          []byte                  `gorm:"column:row_ver"`

	// active version
	AV_ID            mssql.UniqueIdentifier `gorm:"column:av_id"`
	AV_VersionNo     int                    `gorm:"column:av_version_no"`
	AV_PeriodFrom    *time.Time             `gorm:"column:av_period_from"`
	AV_PeriodTo      *time.Time             `gorm:"column:av_period_to"`
	AV_Status        string                 `gorm:"column:av_status"`
	AV_Read          bool                   `gorm:"column:av_read"`
	AV_Note          string                 `gorm:"column:av_note"`
	AV_Quantity      *int                   `gorm:"column:av_quantity"`
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

type ForecastStatusSummary struct {
	VendorCode    string `gorm:"column:vendor_code"`
	NewCount      int64  `gorm:"column:new_count"`
	ConfirmCount  int64  `gorm:"column:confirm_count"`
	RejectCount   int64  `gorm:"column:reject_count"`
	ChangeCount   int64  `gorm:"column:change_count"`
	ApprovedCount int64  `gorm:"column:approved_count"`
	TotalCount    int64  `gorm:"column:total_count"`
}

type ForecastBasicInfo struct {
	NumberForecast string `gorm:"column:number_forecast"`
	VendorCode     string `gorm:"column:vendor_code"`
}
