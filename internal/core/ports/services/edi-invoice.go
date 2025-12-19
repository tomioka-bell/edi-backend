package ports

import (
	"backend/internal/core/domains"
	"backend/internal/core/models"
	"context"

	mssql "github.com/microsoft/go-mssqldb"
)

type EDIInvoiceService interface {
	CreateNewInvoiceWithVersion(ctx context.Context, headerIn *models.EDIInvoiceResp, versionIn *models.EDIInvoiceVersionResp) (*models.EDIInvoiceResp, error)
	MarkInvoiceAsReadService(id mssql.UniqueIdentifier) error
	GetEDIInvoiceWithActiveTopService(limit int, vendorCode string) ([]models.EDIInvoiceWithActiveReq, error)
	GetEDIInvoiceDetailByNumber(number string) (*models.EDIInvoiceDetailResp, error)
	UpdateStatusInvoiceService(id mssql.UniqueIdentifier, status string) error
	CreateEDIInvoiceVersionService(req models.EDIInvoiceVersionResp) error
	GenerateRunningNumberService() (string, error)
	GetInvoiceDetailByNumberOrderService(number string) (*models.EDIInvoiceDetailResp, error)
	GetStatusInvoiceSummaryByVendorCodeService(vendorCode string) (*models.StatusInvoiceSummaryResp, error)

	GetEDIInvoiceVersionByIDService(ediInvoiceVersionID string) (domains.EDIInvoiceVersion, error)
	GetEDIVendorNotificationRecipientByCompanyService(company string) ([]models.EDIVendorNotificationRecipientReq, error)
	GetInvoiceBasicByIDService(ediInvoiceID string) (domains.InvoiceBasicInfo, error)

	// =============================================== Status Log =========================================================
	CreateEDIInvoiceVersionStatusLogService(req models.EDIInvoiceVersionStatusLogResp) error
	GetInvoiceVersionStatusLogByInvoiceVersionIDService(InvoiceVersionID string) ([]models.EDIInvoiceVersionStatusLogReq, error)
	GetInvoiceVersionStatusLogByInvoiceVersionIDAndApprovedService(InvoiceVersionID string) ([]models.EDIInvoiceVersionStatusLogReq, error)
}
