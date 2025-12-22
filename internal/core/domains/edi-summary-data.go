package domains

import "time"

type NumberCountSummary struct {
	NumberValue string `json:"number_value"`
	TotalCount  int    `json:"total_count"`
}

type MonthlyStatusSummary struct {
	Year          int   `json:"year" gorm:"column:year"`
	Month         int   `json:"month" gorm:"column:month"`
	NewCount      int64 `json:"new_count" gorm:"column:new_count"`
	ConfirmCount  int64 `json:"confirm_count" gorm:"column:confirm_count"`
	RejectCount   int64 `json:"reject_count" gorm:"column:reject_count"`
	ChangeCount   int64 `json:"change_count" gorm:"column:change_count"`
	ApprovedCount int64 `json:"approved_count" gorm:"column:approved_count"`
	TotalCount    int64 `json:"total_count" gorm:"column:total_count"`
}

type VendorForecastSummary struct {
	NumberForecast string    `json:"number_forecast"`
	StatusForecast string    `json:"status_forecast"`
	ReadForecast   bool      `json:"read_forecast"`
	VendorCode     string    `json:"vendor_code"`
	CreatedAt      time.Time `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
}

type VendorOrderSummary struct {
	NumberOrder string    `json:"number_order"`
	StatusOrder string    `json:"status_order"`
	ReadOrder   bool      `json:"read_order"`
	VendorCode  string    `json:"vendor_code"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
}

type VendorInvoiceSummary struct {
	NumberInvoice string    `json:"number_invoice"`
	StatusInvoice string    `json:"status_invoice"`
	ReadInvoice   bool      `json:"read_invoice"`
	VendorCode    string    `json:"vendor_code"`
	CreatedAt     time.Time `gorm:"column:created_at;type:datetimeoffset(7);autoCreateTime"`
}

type ForecastPeriodAlert struct {
	NumberForecast  string    `json:"number_forecast"`
	VendorCode      string    `json:"vendor_code"`
	VersionNo       int       `json:"version_no"`
	PeriodFrom      time.Time `json:"period_from"`
	ReadForecast    bool      `json:"read_forecast"`
	TargetTime      time.Time `json:"target_time"`
	SecondsToTarget int64     `json:"seconds_to_target"`
	DaysToTarget    int64     `json:"days_to_target"`
	ReminderDays    int       `json:"reminder_days"`

	Email string `json:"email"`
}

type OrderPeriodAlert struct {
	NumberOrder     string    `json:"number_order"`
	VendorCode      string    `json:"vendor_code"`
	VersionNo       int       `json:"version_no"`
	PeriodFrom      time.Time `json:"period_from"`
	ReadOrder       bool      `json:"read_order"`
	TargetTime      time.Time `json:"target_time"`
	SecondsToTarget int64     `json:"seconds_to_target"`
	DaysToTarget    int64     `json:"days_to_target"`
	ReminderDays    int       `json:"reminder_days"`

	Email string `json:"email"`
}

type InvoicePeriodAlert struct {
	NumberInvoice   string    `json:"number_invoice"`
	VendorCode      string    `json:"vendor_code"`
	VersionNo       int       `json:"version_no"`
	PeriodFrom      time.Time `json:"period_from"`
	ReadInvoice     bool      `json:"read_invoice"`
	TargetTime      time.Time `json:"target_time"`
	SecondsToTarget int64     `json:"seconds_to_target"`
	DaysToTarget    int64     `json:"days_to_target"`
	ReminderDays    int       `json:"reminder_days"`

	Email string `json:"email"`
}
