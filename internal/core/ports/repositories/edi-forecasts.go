package ports

import (
	"backend/internal/core/domains"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

type EDIForecastRepository interface {
	CreateNewForecastWithVersion(db *gorm.DB, header *domains.EDI_Forecast, version *domains.EDI_ForecastVersion) error
	MarkForecastAsRead(id mssql.UniqueIdentifier) (*domains.EDI_Forecast, error)
	CreateEDIForecastRepository(header *domains.EDI_Forecast) error
	CreateEDIForecastVersionRepository(version *domains.EDI_ForecastVersion) error
	UpdateActiveVersion(headerID, versionID string) error
	GetEDIForecastWithActiveTop(limit int, vendorCode string) ([]domains.EDIForecastWithActive, error)
	GetEDIForecastWithAllVersionsTop(limit int) ([]domains.EDI_Forecast, error)
	GetEDIForecastWithActiveByNumber(number string) (*domains.EDIForecastWithActive, error)
	GetForecastHeaderByNumber(number string) (*domains.EDI_Forecast, error)
	GetForecastVersionsByForecastID(forecastID mssql.UniqueIdentifier) ([]domains.EDI_ForecastVersion, error)
	UpdateStatusForecast(id mssql.UniqueIdentifier, status string) error
	GetMaxVersionNoByForecastID(ediForecastID mssql.UniqueIdentifier) (int, error)
	UpdateForecastVersionWithMap(forecastVersionID string, updates map[string]any) error
	UpdateActiveForecastVersion(forecastID mssql.UniqueIdentifier, versionID mssql.UniqueIdentifier) error
	GetEDIForecastVersionByID(ediForecastVersionID string) (domains.EDI_ForecastVersion, error)
	GetEDIVendorNotificationRecipientByCompany(company string) ([]domains.EDIVendorNotificationRecipient, error)
	GetLastForecastRunningForDate(date string) (int, error)
	GetStatusSummaryByVendorCode(vendorCode string) (*domains.ForecastStatusSummary, error)
	GetForecastBasicByID(ediForecastID string) (domains.ForecastBasicInfo, error)

	// =============================================== Status Log =========================================================
	CreateEDIForecastVersionStatusLog(m *domains.EDI_ForecastVersionStatusLog) error
	GetForecastVersionStatusLogByForecastVersionID(forecastVersionID string) ([]domains.EDI_ForecastVersionStatusLog, error)
	GetForecastVersionStatusLogByForecastVersionIDAndApproved(forecastVersionID string) ([]domains.EDI_ForecastVersionStatusLog, error)
}
