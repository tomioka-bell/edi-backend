package models

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

type EDIOrderResp struct {
	EDIOrderID            mssql.UniqueIdentifier  `json:"edi_order_id"`
	NumberOrder           string                  `json:"number_order"`
	VendorCode            string                  `json:"vendor_code"`
	NumberForecast        string                  `json:"number_forecast"`
	ProductCode           *string                 `json:"product_code"`
	ReadOrder             bool                    `json:"read_order"`
	ActiveVersionID       *mssql.UniqueIdentifier `json:"active_version_id"`
	CreatedByExternalID   string                  `json:"created_by_external_id"`
	CreatedBySourceSystem string                  `json:"created_by_source_system"`
	FileURL               *string                 `json:"file_url"`
	CreatedAt             time.Time               `json:"created_at"`
	UpdatedAt             time.Time               `json:"updated_at"`
	DeletedAt             gorm.DeletedAt          `json:"deleted_at"`
	RowVer                []byte                  `json:"row_ver"`

	Versions []EDIOrderVersionResp `json:"versions"`
}

type EDIOrderVersionResp struct {
	EDIOrderVersionID     mssql.UniqueIdentifier `json:"edi_order_version_id"`
	EDIOrderID            mssql.UniqueIdentifier `json:"edi_order_id"`
	VersionNo             int                    `json:"version_no"`
	PeriodFrom            *time.Time             `json:"period_from"`
	PeriodTo              *time.Time             `json:"period_to"`
	StatusOrder           string                 `json:"status_order"`
	ReadOrder             bool                   `json:"read_order"`
	Note                  string                 `json:"note"`
	SourceFileURL         *string                `json:"source_file_url"`
	CreatedByExternalID   string                 `json:"created_by_external_id"`
	CreatedBySourceSystem string                 `json:"created_by_source_system"`
	CreatedAt             time.Time              `json:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
	DeletedAt             gorm.DeletedAt         `json:"deleted_at"`
	RowVer                []byte                 `json:"row_ver"`
}

type EDIOrderWithActiveReq struct {
	// header
	EDIOrderID      mssql.UniqueIdentifier  `json:"edi_order_id"`
	NumberOrder     string                  `json:"number_order"`
	VendorCode      string                  `json:"vendor_code"`
	NumberForecast  string                  `json:"number_forecast"`
	ProductCode     *string                 `json:"product_code"`
	ReadOrder       bool                    `json:"read_order"`
	StatusOrder     string                  `json:"status_order"`
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
	AV_Quantity      *int                   `json:"av_quantity"`
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

type CreateOrderReq struct {
	NumberOrder           string              `json:"number_order"`
	NumberForecast        string              `json:"number_forecast"`
	VendorCode            string              `json:"vendor_code"`
	FileURL               *string             `json:"file_url"`
	CreatedByExternalID   string              `json:"created_by_external_id"`
	CreatedBySourceSystem string              `json:"created_by_source_system"`
	Versions              []CreateOrderVerReq `json:"versions"`
}

type CreateOrderVerReq struct {
	PeriodFrom            *time.Time `json:"period_from"`
	PeriodTo              *time.Time `json:"period_to"`
	StatusOrder           string     `json:"status_order"`
	ReadOrder             bool       `json:"read_order"`
	Note                  string     `json:"note"`
	Quantity              *int       `json:"quantity"`
	SourceFileURL         *string    `json:"source_file_url"`
	CreatedByExternalID   string     `json:"created_by_external_id"`
	CreatedBySourceSystem string     `json:"created_by_source_system"`
}

type EDIOrderDetailResp struct {
	EDIOrderID      mssql.UniqueIdentifier    `json:"edi_order_id"`
	NumberOrder     string                    `json:"number_order"`
	ReadOrder       bool                      `json:"read_order"`
	NumberForecast  string                    `json:"number_forecast"`
	VendorCode      string                    `json:"vendor_code"`
	StatusOrder     string                    `json:"status_order"`
	FileURL         *string                   `json:"file_url"`
	ActiveVersionID *mssql.UniqueIdentifier   `json:"active_version_id"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	Versions        []EDIOrderVersionItemResp `json:"versions"`
}

type EDIOrderVersionItemResp struct {
	EDIOrderVersionID mssql.UniqueIdentifier `json:"edi_order_version_id"`
	VersionNo         int                    `json:"version_no"`
	PeriodFrom        *time.Time             `json:"period_from"`
	PeriodTo          *time.Time             `json:"period_to"`
	StatusOrder       string                 `json:"status_order"`
	Note              string                 `json:"note"`
	Quantity          *int                   `json:"quantity"`
	SourceFileURL     *string                `json:"source_file_url"`
	CreatedAt         time.Time              `json:"created_at"`
	IsActive          bool                   `json:"is_active"`

	CreatedBy *CreatedByPrincipal `json:"created_by,omitempty"`
}

type StatusOrderSummaryResp struct {
	VendorCode    string `json:"vendor_code"`
	NewCount      int64  `json:"new_count"`
	ConfirmCount  int64  `json:"confirm_count"`
	RejectCount   int64  `json:"reject_count"`
	ApprovedCount int64  `json:"approved_count"`
	TotalCount    int64  `json:"total_count"`
}

type EDIOrderByNumberForecastResp struct {
	EDIOrderID     mssql.UniqueIdentifier `json:"edi_order_id"`
	NumberOrder    string                 `json:"number_order"`
	NumberForecast string                 `json:"number_forecast"`
	StatusOrder    string                 `json:"status_order"`
	VendorCode     string                 `json:"vendor_code"`
	CreatedAt      time.Time              `json:"created_at"`
	PeriodTo       *time.Time             `json:"period_to"`
}
