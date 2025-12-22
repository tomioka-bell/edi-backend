package services

import (
	"context"

	"backend/internal/core/models"
	ports "backend/internal/core/ports/repositories"
	servicesports "backend/internal/core/ports/services"
)

type EDISummaryDataService struct {
	EDISummaryDataRepo ports.SummaryDataRepository
}

func NewEDISummaryDataService(EDISummaryDataRepo ports.SummaryDataRepository) servicesports.EDISummaryDataService {
	return &EDISummaryDataService{EDISummaryDataRepo: EDISummaryDataRepo}
}

func (s *EDISummaryDataService) GetAllStatusSummaryData() (*models.AllSummary, error) {

	forecast, err := s.EDISummaryDataRepo.GetStatusForecastSummaryByVendorCode()
	if err != nil {
		return nil, err
	}

	order, err := s.EDISummaryDataRepo.GetStatusOrderSummaryByVendorCode()
	if err != nil {
		return nil, err
	}

	invoice, err := s.EDISummaryDataRepo.GetStatusInvoiceSummaryByVendorCode()
	if err != nil {
		return nil, err
	}

	return &models.AllSummary{
		Forecast: forecast,
		Order:    order,
		Invoice:  invoice,
	}, nil
}

func (s *EDISummaryDataService) GetAllTotalCountSummary() (*models.TotalCountSummary, error) {
	forecast, err := s.EDISummaryDataRepo.GetForecastTotalCount()
	if err != nil {
		return nil, err
	}

	order, err := s.EDISummaryDataRepo.GetOrderTotalCount()
	if err != nil {
		return nil, err
	}

	invoice, err := s.EDISummaryDataRepo.GetInvoiceTotalCount()
	if err != nil {
		return nil, err
	}

	return &models.TotalCountSummary{
		Forecast: forecast,
		Order:    order,
		Invoice:  invoice,
	}, nil
}

func (s *EDISummaryDataService) GetAllStatusTotalSummary() (*models.AllStatusTotalSummary, error) {
	forecast, err := s.EDISummaryDataRepo.GetForecastStatusTotal()
	if err != nil {
		return nil, err
	}

	order, err := s.EDISummaryDataRepo.GetOrderStatusTotal()
	if err != nil {
		return nil, err
	}

	invoice, err := s.EDISummaryDataRepo.GetInvoiceStatusTotal()
	if err != nil {
		return nil, err
	}

	return &models.AllStatusTotalSummary{
		Forecast: forecast,
		Order:    order,
		Invoice:  invoice,
	}, nil
}

func (s *EDISummaryDataService) GetAllMonthlyStatusSummary() (*models.AllMonthlyStatusSummary, error) {
	forecast, err := s.EDISummaryDataRepo.GetForecastMonthlyStatusSummary()
	if err != nil {
		return nil, err
	}

	order, err := s.EDISummaryDataRepo.GetOrderMonthlyStatusSummary()
	if err != nil {
		return nil, err
	}

	invoice, err := s.EDISummaryDataRepo.GetInvoiceMonthlyStatusSummary()
	if err != nil {
		return nil, err
	}

	return &models.AllMonthlyStatusSummary{
		Forecast: forecast,
		Order:    order,
		Invoice:  invoice,
	}, nil
}

func (s *EDISummaryDataService) CountUserService() (int64, error) {
	return s.EDISummaryDataRepo.CountUsers()
}

func (s *EDISummaryDataService) GetVendorFlatSummary(ctx context.Context, vendorCode string) ([]map[string]any, error) {

	var response []map[string]any

	// Forecast
	forecasts, err := s.EDISummaryDataRepo.GetUnreadForecastByVendor(ctx, vendorCode)
	if err != nil {
		return nil, err
	}
	for _, f := range forecasts {
		response = append(response, map[string]any{
			"number_forecast": f.NumberForecast,
			"status_forecast": f.StatusForecast,
			"vendor_code":     f.VendorCode,
			"created_at":      f.CreatedAt,
		})
	}

	// Order
	orders, err := s.EDISummaryDataRepo.GetUnreadOrderByVendor(ctx, vendorCode)
	if err != nil {
		return nil, err
	}
	for _, o := range orders {
		response = append(response, map[string]any{
			"number_order": o.NumberOrder,
			"status_order": o.StatusOrder,
			"vendor_code":  o.VendorCode,
			"created_at":   o.CreatedAt,
		})
	}

	if vendorCode == "Prospira (Thailand) Co., Ltd." {
		invoices, err := s.EDISummaryDataRepo.GetUnreadInvoiceByVendor(ctx, vendorCode)
		if err != nil {
			return nil, err
		}
		for _, i := range invoices {
			response = append(response, map[string]any{
				"number_invoice": i.NumberInvoice,
				"status_invoice": i.StatusInvoice,
				"vendor_code":    i.VendorCode,
				"created_at":     i.CreatedAt,
			})
		}
	}

	return response, nil
}

func (s *EDISummaryDataService) GetForecastPeriodAlerts(ctx context.Context) ([]models.ForecastPeriodAlertResponse, error) {
	rows, err := s.EDISummaryDataRepo.GetForecastPeriodAlerts(ctx)
	if err != nil {
		return nil, err
	}

	type key struct {
		VendorCode     string
		NumberForecast string
		VersionNo      int
	}

	m := make(map[key]*models.ForecastPeriodAlertResponse)

	for _, r := range rows {
		k := key{VendorCode: r.VendorCode, NumberForecast: r.NumberForecast, VersionNo: r.VersionNo}

		if _, ok := m[k]; !ok {
			m[k] = &models.ForecastPeriodAlertResponse{
				NumberForecast:  r.NumberForecast,
				VendorCode:      r.VendorCode,
				VersionNo:       r.VersionNo,
				PeriodFrom:      r.PeriodFrom,
				ReadForecast:    r.ReadForecast,
				TargetTime:      r.TargetTime,
				SecondsToTarget: r.SecondsToTarget,
				DaysToTarget:    r.DaysToTarget,
				ReminderDays:    r.ReminderDays,
				Emails:          []string{},
			}
		}

		if r.Email != "" {
			m[k].Emails = append(m[k].Emails, r.Email)
		}
	}

	resp := make([]models.ForecastPeriodAlertResponse, 0, len(m))
	for _, v := range m {
		resp = append(resp, *v)
	}

	return resp, nil
}

func (s *EDISummaryDataService) GetOrderPeriodAlerts(ctx context.Context) ([]models.OrderPeriodAlertResponse, error) {
	rows, err := s.EDISummaryDataRepo.GetOrderPeriodAlerts(ctx)
	if err != nil {
		return nil, err
	}

	type key struct {
		VendorCode  string
		NumberOrder string
		VersionNo   int
	}

	m := make(map[key]*models.OrderPeriodAlertResponse)

	for _, r := range rows {
		k := key{VendorCode: r.VendorCode, NumberOrder: r.NumberOrder, VersionNo: r.VersionNo}

		if _, ok := m[k]; !ok {
			m[k] = &models.OrderPeriodAlertResponse{
				NumberOrder:     r.NumberOrder,
				VendorCode:      r.VendorCode,
				VersionNo:       r.VersionNo,
				PeriodFrom:      r.PeriodFrom,
				ReadOrder:       r.ReadOrder,
				TargetTime:      r.TargetTime,
				SecondsToTarget: r.SecondsToTarget,
				DaysToTarget:    r.DaysToTarget,
				ReminderDays:    r.ReminderDays,
				Emails:          []string{},
			}
		}

		if r.Email != "" {
			m[k].Emails = append(m[k].Emails, r.Email)
		}
	}

	resp := make([]models.OrderPeriodAlertResponse, 0, len(m))
	for _, v := range m {
		resp = append(resp, *v)
	}

	return resp, nil
}

func (s *EDISummaryDataService) GetInvoicePeriodAlerts(ctx context.Context) ([]models.InvoicePeriodAlertResponse, error) {
	rows, err := s.EDISummaryDataRepo.GetInvoicePeriodAlerts(ctx)
	if err != nil {
		return nil, err
	}

	type key struct {
		VendorCode    string
		NumberInvoice string
		VersionNo     int
	}

	m := make(map[key]*models.InvoicePeriodAlertResponse)

	for _, r := range rows {
		k := key{VendorCode: r.VendorCode, NumberInvoice: r.NumberInvoice, VersionNo: r.VersionNo}

		if _, ok := m[k]; !ok {
			m[k] = &models.InvoicePeriodAlertResponse{
				NumberInvoice:   r.NumberInvoice,
				VendorCode:      r.VendorCode,
				VersionNo:       r.VersionNo,
				PeriodFrom:      r.PeriodFrom,
				ReadInvoice:     r.ReadInvoice,
				TargetTime:      r.TargetTime,
				SecondsToTarget: r.SecondsToTarget,
				DaysToTarget:    r.DaysToTarget,
				ReminderDays:    r.ReminderDays,
				Emails:          []string{},
			}
		}

		if r.Email != "" {
			m[k].Emails = append(m[k].Emails, r.Email)
		}
	}

	resp := make([]models.InvoicePeriodAlertResponse, 0, len(m))
	for _, v := range m {
		resp = append(resp, *v)
	}

	return resp, nil
}
