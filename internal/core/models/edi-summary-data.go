package models

import (
	"backend/internal/core/domains"
	"time"
)

type AllSummary struct {
	Forecast []domains.ForecastStatusSummary `json:"forecast"`
	Order    []domains.OrderStatusSummary    `json:"order"`
	Invoice  []domains.OrderStatusSummary    `json:"invoice"`
}

type TotalCountSummary struct {
	Forecast int64 `json:"forecast"`
	Order    int64 `json:"order"`
	Invoice  int64 `json:"invoice"`
}

type AllStatusTotalSummary struct {
	Forecast *domains.StatusTotalSummary `json:"forecast"`
	Order    *domains.StatusTotalSummary `json:"order"`
	Invoice  *domains.StatusTotalSummary `json:"invoice"`
}

type AllMonthlyStatusSummary struct {
	Forecast []domains.MonthlyStatusSummary `json:"forecast"`
	Order    []domains.MonthlyStatusSummary `json:"order"`
	Invoice  []domains.MonthlyStatusSummary `json:"invoice"`
}

type ForecastPeriodAlertResponse struct {
	NumberForecast  string    `json:"number_forecast"`
	VendorCode      string    `json:"vendor_code"`
	VersionNo       int       `json:"version_no"`
	PeriodFrom      time.Time `json:"period_from"`
	ReadForecast    bool      `json:"read_forecast"`
	TargetTime      time.Time `json:"target_time"`
	SecondsToTarget int64     `json:"seconds_to_target"`
	DaysToTarget    int64     `json:"days_to_target"`
	ReminderDays    int       `json:"reminder_days"`
	Emails          []string  `json:"emails"`
}

type OrderPeriodAlertResponse struct {
	NumberOrder     string    `json:"number_order"`
	VendorCode      string    `json:"vendor_code"`
	VersionNo       int       `json:"version_no"`
	PeriodFrom      time.Time `json:"period_from"`
	ReadOrder       bool      `json:"read_order"`
	TargetTime      time.Time `json:"target_time"`
	SecondsToTarget int64     `json:"seconds_to_target"`
	DaysToTarget    int64     `json:"days_to_target"`
	ReminderDays    int       `json:"reminder_days"`
	Emails          []string  `json:"emails"`
}

type InvoicePeriodAlertResponse struct {
	NumberInvoice   string    `json:"number_invoice"`
	VendorCode      string    `json:"vendor_code"`
	VersionNo       int       `json:"version_no"`
	PeriodFrom      time.Time `json:"period_from"`
	ReadInvoice     bool      `json:"read_invoice"`
	TargetTime      time.Time `json:"target_time"`
	SecondsToTarget int64     `json:"seconds_to_target"`
	DaysToTarget    int64     `json:"days_to_target"`
	ReminderDays    int       `json:"reminder_days"`
	Emails          []string  `json:"emails"`
}
