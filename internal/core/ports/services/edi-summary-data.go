package ports

import (
	"backend/internal/core/models"
	"context"
)

type EDISummaryDataService interface {
	GetAllStatusSummaryData() (*models.AllSummary, error)
	GetAllTotalCountSummary() (*models.TotalCountSummary, error)
	GetAllStatusTotalSummary() (*models.AllStatusTotalSummary, error)
	GetAllMonthlyStatusSummary() (*models.AllMonthlyStatusSummary, error)
	CountUserService() (int64, error)
	GetVendorFlatSummary(ctx context.Context, vendorCode string) ([]map[string]any, error)
	GetForecastPeriodAlerts(ctx context.Context) ([]models.ForecastPeriodAlertResponse, error)
	GetOrderPeriodAlerts(ctx context.Context) ([]models.OrderPeriodAlertResponse, error)
	GetInvoicePeriodAlerts(ctx context.Context) ([]models.InvoicePeriodAlertResponse, error)
}
