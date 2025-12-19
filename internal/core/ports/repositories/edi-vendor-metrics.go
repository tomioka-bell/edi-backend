package ports

import "backend/internal/core/domains"

type EDIVendorMetricsRepository interface {
	CreateVendorMetrics(vendorMetrics *domains.EDIVendorMetrics) error
	GetVendorMetricsByCompany(company string) ([]domains.EDIVendorMetrics, error)
	GetAllEDIVendorMetricsTop(limit int) ([]domains.EDIVendorMetrics, error)
	UpdateVendorMetricsWithMap(vendorMetricsID string, updates map[string]interface{}) error
	DeleteVendorMetrics(vendorMetricsID string) error
	GetAllVendorMetrics() ([]domains.EDIVendorMetrics, error)
}
