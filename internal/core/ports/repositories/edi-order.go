package ports

import (
	mssql "github.com/microsoft/go-mssqldb"

	"backend/internal/core/domains"
)

type EDIOrderRepository interface {
	CreateEDIOrderRepository(m *domains.EDIOrder) error
	CreateEDIOrderVersionRepository(version *domains.EDIOrderVersion) error
	GetEDIOrderWithActiveTop(limit int, vendorCode string) ([]domains.EDIOrderWithActive, error)
	GetLastOrderRunningForDate(date string) (int, error)
	UpdateActiveOrderVersion(headerID string, activeVersionID string) error
	MarkOrderAsRead(id mssql.UniqueIdentifier) error
	UpdateStatusOrder(id mssql.UniqueIdentifier, status string) error
	GetMaxVersionNoByOrderID(orderID mssql.UniqueIdentifier) (int, error)
	GetEDIOrderByNumberOrderData(orderNumber string) ([]domains.EDIOrderVersionStatusLog, error)

	GetOrderVersionsByOrderID(orderID mssql.UniqueIdentifier) ([]domains.EDIOrderVersion, error)
	GetOrderHeaderByNumber(number string) (*domains.EDIOrder, error)

	GetOrderVersionStatusLogByOrderVersionID(orderVersionID string) ([]domains.EDIOrderVersionStatusLog, error)
	GetOrderVersionStatusLogByOrderVersionIDAndApproved(orderVersionID string) ([]domains.EDIOrderVersionStatusLog, error)
	CreateEDIOrderVersionStatusLog(log *domains.EDIOrderVersionStatusLog) error

	GetStatusOrderSummaryByVendorCode(vendorCode string) (*domains.OrderStatusSummary, error)
	GetOrderHeaderByVendorCode(VendorCode string) (*domains.EDIOrder, error)

	GetOrderByNumberForecast(numberForecast string) ([]domains.EDIOrderHeaderWithPeriod, error)

	GetEDIOrderVersionByID(ediOrderVersionID string) (domains.EDIOrderVersion, error)

	GetEDIVendorNotificationRecipientByCompany(company string) ([]domains.EDIVendorNotificationRecipient, error)
	GetOrderBasicByID(ediOrderID string) (domains.OrderBasicInfo, error)
	GetOrderVersionStatusLogByOrderNumberAndApproved(orderNumber string) ([]domains.EDIOrderVersionStatusLog, error)
}
