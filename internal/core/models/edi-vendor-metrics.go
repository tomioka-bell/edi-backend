package models

import (
	mssql "github.com/microsoft/go-mssqldb"
)

type VendorMetricsResp struct {
	VendorMetricsID mssql.UniqueIdentifier `json:"vendor_metrics_id"`
	Initials        string                 `json:"initials"`
	CompanyName     string                 `json:"company_name"`

	ReminderDays int  `json:"reminder_days"`
	Active       bool `json:"active"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type VendorMetricsReq struct {
	VendorMetricsID mssql.UniqueIdentifier `json:"vendor_metrics_id"`
	Initials        string                 `json:"initials"`
	CompanyName     string                 `json:"company_name"`
	ReminderDays    int                    `json:"reminder_days"`
	Active          bool                   `json:"active"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type VendorMetricsUpdate struct {
	Initials     *string `json:"initials"`
	CompanyName  *string `json:"company_name"`
	ReminderDays int     `json:"reminder_days"`
	Active       *bool   `json:"active"`
}

type VendorMetricsCompanyResp struct {
	Initials    string `json:"initials"`
	CompanyName string `json:"company_name"`
}
