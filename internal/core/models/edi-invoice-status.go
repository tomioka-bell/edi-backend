package models

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
)

type EDIInvoiceVersionStatusLogResp struct {
	EDIInvoiceVersionID   mssql.UniqueIdentifier `json:"edi_invoice_version_id"`
	EDIInvoiceID          mssql.UniqueIdentifier `json:"edi_invoice_id"`
	OldStatus             string                 `json:"old_status"`
	NewStatus             string                 `json:"new_status"`
	Note                  string                 `json:"note"`
	ChangedByExternalID   string                 `json:"created_by_external_id"`
	ChangedBySourceSystem string                 `json:"created_by_source_system"`
	FileURL               *string                `json:"file_url"`
	CreatedAt             time.Time              `json:"created_at"`
}

// ใช้รับข้อมูลจาก multipart/form-data
type EDIInvoiceVersionStatusLogForm struct {
	EDIInvoiceVersionID   string `form:"edi_invoice_version_id"`
	EDIInvoiceID          string `form:"edi_invoice_id"`
	OldStatus             string `form:"old_status"`
	NewStatus             string `form:"new_status"`
	Note                  string `form:"note"`
	ChangedByExternalID   string `form:"created_by_external_id"`
	ChangedBySourceSystem string `form:"created_by_source_system"`
}

type EDIInvoiceVersionStatusLogReq struct {
	EDIInvoiceID          mssql.UniqueIdentifier `json:"edi_invoice_id"`
	OldStatus             string                 `json:"old_status"`
	NewStatus             string                 `json:"new_status"`
	Note                  string                 `json:"note"`
	ChangedByExternalID   string                 `json:"created_by_external_id"`
	ChangedBySourceSystem string                 `json:"created_by_source_system"`
	FileURL               *string                `json:"file_url"`
	CreatedAt             time.Time              `json:"created_at"`

	ChangedByUser *PrincipalResp `json:"changed_by_user,omitempty"`
}
