package ports

import (
	"backend/internal/core/domains"
	"backend/internal/core/models"
	"context"

	mssql "github.com/microsoft/go-mssqldb"
)

type EDIOrderService interface {
	CreateNewOrderWithVersion(ctx context.Context, headerIn *models.EDIOrderResp, versionIn *models.EDIOrderVersionResp) (*models.EDIOrderResp, error)
	MarkOrderAsReadService(id mssql.UniqueIdentifier) error
	GetEDIOrderWithActiveTopService(limit int, vendorCode string) ([]models.EDIOrderWithActiveReq, error)
	GetEDIOrderDetailByNumber(number string) (*models.EDIOrderDetailResp, error)
	UpdateStatusOrderService(id mssql.UniqueIdentifier, status string) error
	CreateEDIOrderVersionService(req models.EDIOrderVersionResp) error
	GetEDIOrderByNumberOrderDataService(numberOrder string) ([]models.EDIOrderVersionStatusLogReq, error)
	GenerateRunningNumberService() (string, error)
	GetStatusOrderSummaryByVendorCodeService(vendorCode string) (*models.StatusOrderSummaryResp, error)
	GetOrderHeaderByVendorCodeService(VendorCode string) ([]models.EDIOrderByvendorReq, error)
	GetOrderHeaderByNumberForecastService(numberForecast string) ([]models.EDIOrderByNumberForecastResp, error)
	GetEDIOrderVersionByIDService(ediOrderVersionID string) (domains.EDIOrderVersion, error)
	GetEDIVendorNotificationRecipientByCompanyService(company string) ([]models.EDIVendorNotificationRecipientReq, error)
	GetOrderBasicByIDService(ediOrderID string) (domains.OrderBasicInfo, error)

	// =============================================== Status Log =========================================================
	CreateEDIOrderVersionStatusLogService(req models.EDIOrderVersionStatusLogResp) error
	GetOrderVersionStatusLogByOrderVersionIDService(orderVersionID string) ([]models.EDIOrderVersionStatusLogReq, error)
	GetOrderVersionStatusLogByOrderVersionIDAndApprovedService(orderVersionID string) ([]models.EDIOrderVersionStatusLogReq, error)
}
