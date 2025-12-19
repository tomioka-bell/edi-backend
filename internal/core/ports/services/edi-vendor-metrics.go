package ports

import (
	"backend/internal/core/models"

	mssql "github.com/microsoft/go-mssqldb"
)

type EDIVendorMetricsService interface {
	CreateVendorMetricsService(req models.VendorMetricsResp) error
	GetVendorMetricsByCompanyService(company string) ([]models.VendorMetricsResp, error)
	GetAllEDIVendorMetricsTopService(limit int) ([]models.VendorMetricsResp, error)
	UpdateVendorMetricsService(id mssql.UniqueIdentifier, updates map[string]interface{}) error
	DeleteVendorMetricsService(vendorMetricsID mssql.UniqueIdentifier) error
	GetAllVendorMetricservice() ([]models.VendorMetricsCompanyResp, error)
}
