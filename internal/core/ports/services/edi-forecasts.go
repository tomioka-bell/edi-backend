package ports

import (
	"backend/internal/core/domains"
	"backend/internal/core/models"
	"context"

	mssql "github.com/microsoft/go-mssqldb"
)

type EDIForecastService interface {
	CreateNewForecastWithVersion(ctx context.Context, headerIn *models.EDI_ForecastResp, versionIn *models.EDI_ForecastVersionResp) (*models.EDI_ForecastResp, error)
	GetEDIForecastWithActiveTopService(limit int, vendorCode string) ([]models.EDIForecastWithActiveReq, error)
	GetEDIForecastWithActiveByNumberService(number string) (*models.EDIForecastWithActiveReq, error)
	CreateEDIForecastVersionService(req models.EDI_ForecastVersionResp) error
	GetEDIForecastDetailByNumber(number string) (*models.EDIForecastDetailResp, error)
	MarkForecastAsReadService(id mssql.UniqueIdentifier) (*domains.EDI_Forecast, error)
	UpdateStatusForecastService(id mssql.UniqueIdentifier, status string) error
	GetEDIForecastVersionByIDService(ediForecastVersionID string) (domains.EDI_ForecastVersion, error)
	GetEDIVendorNotificationRecipientByCompanyService(company string) ([]models.EDIVendorNotificationRecipientReq, error)
	GenerateRunningNumberService() (string, error)
	GetStatusSummaryByVendorCodeService(vendorCode string) (*models.ForecastStatusSummaryResp, error)
	GetForecastBasicByIDService(ediForecastID string) (domains.ForecastBasicInfo, error)
	// =============================================== Status Log =========================================================
	CreateEDIForecastVersionStatusLogService(req models.EDI_ForecastVersionStatusLogResp) error
	GetForecastVersionStatusLogByForecastVersionIDService(forecastVersionID string) ([]models.EDI_ForecastVersionStatusLogReq, error)
	GetForecastVersionStatusLogByForecastVersionIDAndApprovedService(forecastVersionID string) ([]models.EDI_ForecastVersionStatusLogReq, error)
}
