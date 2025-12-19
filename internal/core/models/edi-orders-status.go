package models

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
)

type EDIOrderVersionStatusLogResp struct {
	EDIOrderVersionID     mssql.UniqueIdentifier `json:"edi_order_version_id"`
	EDIOrderID            mssql.UniqueIdentifier `json:"edi_order_id"`
	OldStatus             string                 `json:"old_status"`
	NewStatus             string                 `json:"new_status"`
	Note                  string                 `json:"note"`
	ChangedByExternalID   string                 `json:"created_by_external_id"`
	ChangedBySourceSystem string                 `json:"created_by_source_system"`
	FileURL               *string                `json:"file_url"`
	CreatedAt             time.Time              `json:"created_at"`
}

// ใช้รับข้อมูลจาก multipart/form-data
type EDIOrderVersionStatusLogForm struct {
	EDIOrderVersionID     string `form:"edi_order_version_id"`
	EDIOrderID            string `form:"edi_order_id"`
	OldStatus             string `form:"old_status"`
	NewStatus             string `form:"new_status"`
	Note                  string `form:"note"`
	ChangedByExternalID   string `form:"created_by_external_id"`
	ChangedBySourceSystem string `form:"created_by_source_system"`
}

type EDIOrderVersionStatusLogReq struct {
	EDIOrderID            mssql.UniqueIdentifier `json:"edi_order_id"`
	OldStatus             string                 `json:"old_status"`
	NewStatus             string                 `json:"new_status"`
	Note                  string                 `json:"note"`
	ChangedByExternalID   string                 `json:"created_by_external_id"`
	ChangedBySourceSystem string                 `json:"created_by_source_system"`
	FileURL               *string                `json:"file_url"`
	CreatedAt             time.Time              `json:"created_at"`

	ChangedByUser *PrincipalResp `json:"changed_by_user,omitempty"`
}

type EDIOrderByvendorReq struct {
	NumberForecast string `json:"number_forecast"`
}
