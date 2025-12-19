package ports

import (
	"backend/internal/core/domains"

	mssql "github.com/microsoft/go-mssqldb"
)

type EDIInvoiceRepository interface {
	CreateEDIInvoiceRepository(m *domains.EDIInvoice) error
	CreateEDIInvoiceVersionRepository(version *domains.EDIInvoiceVersion) error
	GetEDIInvoiceWithActiveTop(limit int, vendorCode string) ([]domains.EDIInvoiceWithActive, error)
	UpdateActiveInvoiceVersion(headerID string, activeVersionID string) error
	MarkInvoiceAsRead(id mssql.UniqueIdentifier) error
	UpdateStatusInvoice(id mssql.UniqueIdentifier, status string) error
	GetMaxVersionNoByInvoiceID(InvoiceID mssql.UniqueIdentifier) (int, error)
	GetLastInvoiceRunningForDate(date string) (int, error)
	GetInvoiceByNumberOrder(order string) (*domains.EDIInvoice, error)

	GetInvoiceVersionsByInvoiceID(InvoiceID mssql.UniqueIdentifier) ([]domains.EDIInvoiceVersion, error)
	GetInvoiceHeaderByNumber(number string) (*domains.EDIInvoice, error)

	GetInvoiceBasicByID(ediInvoiceID string) (domains.InvoiceBasicInfo, error)

	GetInvoiceVersionStatusLogByInvoiceVersionID(InvoiceVersionID string) ([]domains.EDIInvoiceVersionStatusLog, error)
	GetInvoiceVersionStatusLogByInvoiceVersionIDAndApproved(InvoiceVersionID string) ([]domains.EDIInvoiceVersionStatusLog, error)
	CreateEDIInvoiceVersionStatusLog(log *domains.EDIInvoiceVersionStatusLog) error

	GetStatusInvoiceSummaryByVendorCode(vendorCode string) (*domains.InvoiceStatusSummary, error)

	GetEDIInvoiceVersionByID(ediInvoiceVersionID string) (domains.EDIInvoiceVersion, error)
	GetEDIVendorNotificationRecipientByCompany(company string) ([]domains.EDIVendorNotificationRecipient, error)
}
