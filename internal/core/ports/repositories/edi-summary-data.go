package ports

import (
	"backend/internal/core/domains"
	"context"
)

type SummaryDataRepository interface {
	GetStatusForecastSummaryByVendorCode() ([]domains.ForecastStatusSummary, error)
	GetStatusOrderSummaryByVendorCode() ([]domains.OrderStatusSummary, error)
	GetStatusInvoiceSummaryByVendorCode() ([]domains.OrderStatusSummary, error)

	GetForecastTotalCount() (int64, error)
	GetOrderTotalCount() (int64, error)
	GetInvoiceTotalCount() (int64, error)

	GetForecastStatusTotal() (*domains.StatusTotalSummary, error)
	GetOrderStatusTotal() (*domains.StatusTotalSummary, error)
	GetInvoiceStatusTotal() (*domains.StatusTotalSummary, error)

	GetForecastMonthlyStatusSummary() ([]domains.MonthlyStatusSummary, error)
	GetOrderMonthlyStatusSummary() ([]domains.MonthlyStatusSummary, error)
	GetInvoiceMonthlyStatusSummary() ([]domains.MonthlyStatusSummary, error)

	CountUsers() (int64, error)

	GetUnreadForecastByVendor(ctx context.Context, vendorCode string) ([]domains.VendorForecastSummary, error)
	GetUnreadOrderByVendor(ctx context.Context, vendorCode string) ([]domains.VendorOrderSummary, error)
	GetUnreadInvoiceByVendor(ctx context.Context, vendorCode string) ([]domains.VendorInvoiceSummary, error)

	GetForecastPeriodAlerts(ctx context.Context) ([]domains.ForecastPeriodAlert, error)
	GetOrderPeriodAlerts(ctx context.Context) ([]domains.OrderPeriodAlert, error)
	GetInvoicePeriodAlerts(ctx context.Context) ([]domains.InvoicePeriodAlert, error)
}
