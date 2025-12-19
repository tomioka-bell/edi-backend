package models

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

type EDI_ForecastResp struct {
	EDI_ForecastID        mssql.UniqueIdentifier  `json:"edi_forecast_id"`
	NumberForecast        string                  `json:"number_forecast"`
	VendorCode            string                  `json:"vendor_code"`
	ReadForecast          bool                    `json:"read_forecast"`
	ActiveVersionID       *mssql.UniqueIdentifier `json:"active_version_id"`
	CreatedByExternalID   string                  `json:"created_by_external_id"`
	CreatedBySourceSystem string                  `json:"created_by_source_system"`
	FileURL               *string                 `json:"file_url"`
	CreatedAt             time.Time               `json:"created_at"`
	UpdatedAt             time.Time               `json:"updated_at"`
	DeletedAt             gorm.DeletedAt          `json:"deleted_at"`
	RowVer                []byte                  `json:"row_ver"`

	Versions []EDI_ForecastVersionResp `json:"versions"`
}

type EDI_ForecastVersionResp struct {
	EDIForecastVersionID  mssql.UniqueIdentifier `json:"edi_forecast_version_id"`
	EDIForecastID         mssql.UniqueIdentifier `json:"edi_forecast_id"`
	VersionNo             int                    `json:"version_no"`
	PeriodFrom            *time.Time             `json:"period_from"`
	PeriodTo              *time.Time             `json:"period_to"`
	StatusForecast        string                 `json:"status_forecast"`
	ReadForecast          bool                   `json:"read_forecast"`
	Note                  string                 `json:"note"`
	SourceFileURL         *string                `json:"source_file_url"`
	CreatedByExternalID   string                 `json:"created_by_external_id"`
	CreatedBySourceSystem string                 `json:"created_by_source_system"`
	CreatedAt             time.Time              `json:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
	DeletedAt             gorm.DeletedAt         `json:"deleted_at"`
	RowVer                []byte                 `json:"row_ver"`
}

type CreateForecastReq struct {
	NumberForecast        string                 `json:"number_forecast"`
	VendorCode            string                 `json:"vendor_code"`
	FileURL               *string                `json:"file_url"`
	CreatedByExternalID   string                 `json:"created_by_external_id"`
	CreatedBySourceSystem string                 `json:"created_by_source_system"`
	Versions              []CreateForecastVerReq `json:"versions"`
}

type CreateForecastVerReq struct {
	PeriodFrom            *time.Time `json:"period_from"`
	PeriodTo              *time.Time `json:"period_to"`
	StatusForecast        string     `json:"status_forecast"`
	ReadForecast          bool       `json:"read_forecast"`
	Note                  string     `json:"note"`
	SourceFileURL         *string    `json:"source_file_url"`
	CreatedByExternalID   string     `json:"created_by_external_id"`
	CreatedBySourceSystem string     `json:"created_by_source_system"`
}

type EDIForecastWithActiveReq struct {
	// header
	EDIForecastID   mssql.UniqueIdentifier  `json:"edi_forecast_id"`
	NumberForecast  string                  `json:"number_forecast"`
	VendorCode      string                  `json:"vendor_code"`
	ReadForecast    bool                    `json:"read_forecast"`
	StatusForecast  string                  `json:"status_forecast"`
	FileURL         *string                 `json:"file_url"`
	ActiveVersionID *mssql.UniqueIdentifier `json:"active_version_id"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`

	// active version
	AV_ID            mssql.UniqueIdentifier `json:"av_id"`
	AV_VersionNo     int                    `json:"av_version_no"`
	AV_PeriodFrom    *time.Time             `json:"av_period_from"`
	AV_PeriodTo      *time.Time             `json:"av_period_to"`
	AV_Status        string                 `json:"av_status"`
	AV_Read          bool                   `json:"av_read"`
	AV_Note          string                 `json:"av_note"`
	AV_SourceFileURL *string                `json:"av_source_file_url"`
	AV_CreatedAt     time.Time              `json:"av_created_at"`

	// latest status log of active version
	LastStatusLogID *mssql.UniqueIdentifier `json:"last_status_log_id,omitempty"`
	LastOldStatus   *string                 `json:"last_old_status,omitempty"`
	LastNewStatus   *string                 `json:"last_new_status,omitempty"`
	LastStatusNote  *string                 `json:"last_status_note,omitempty"`
	LastFileURL     *string                 `json:"last_file_url,omitempty"`
	LastStatusAt    *time.Time              `json:"last_status_at,omitempty"`
}

type EDIForecastDetailResp struct {
	EDI_ForecastID  mssql.UniqueIdentifier       `json:"edi_forecast_id"`
	NumberForecast  string                       `json:"number_forecast"`
	VendorCode      string                       `json:"vendor_code"`
	StatusForecast  string                       `json:"status_forecast"`
	FileURL         *string                      `json:"file_url"`
	ReadForecast    bool                         `json:"read_forecast"`
	ActiveVersionID *mssql.UniqueIdentifier      `json:"active_version_id"`
	CreatedAt       time.Time                    `json:"created_at"`
	UpdatedAt       time.Time                    `json:"updated_at"`
	Versions        []EDIForecastVersionItemResp `json:"versions"`
}

type CreatedByPrincipal struct {
	ExternalID   string `json:"user_id"`
	Email        string `json:"email"`
	Display_name string `json:"display_name"`
	Profile      string `json:"profile"`
	Group        string `json:"group"`
	Role         string `json:"role"`
	SourceSystem string `json:"source_system"`
	Status       string `json:"status"`
	Username     string `json:"username"`
}

type EDIForecastVersionItemResp struct {
	EDIForecastVersionID mssql.UniqueIdentifier `json:"edi_forecast_version_id"`
	VersionNo            int                    `json:"version_no"`
	PeriodFrom           *time.Time             `json:"period_from"`
	PeriodTo             *time.Time             `json:"period_to"`
	StatusForecast       string                 `json:"status_forecast"`
	ReadForecast         bool                   `json:"read_forecast"`
	Note                 string                 `json:"note"`
	SourceFileURL        *string                `json:"source_file_url"`
	CreatedAt            time.Time              `json:"created_at"`
	IsActive             bool                   `json:"is_active"`

	CreatedBy *CreatedByPrincipal `json:"created_by,omitempty"`
}

type ForecastStatusSummaryResp struct {
	VendorCode    string `json:"vendor_code"`
	NewCount      int64  `json:"new_count"`
	ConfirmCount  int64  `json:"confirm_count"`
	RejectCount   int64  `json:"reject_count"`
	ApprovedCount int64  `json:"approved_count"`
	TotalCount    int64  `json:"total_count"`
}
