package models

import (
	mssql "github.com/microsoft/go-mssqldb"
)

type EDICompanyResp struct {
	CompanyID mssql.UniqueIdentifier `json:"company_id"`
	Name      string                 `json:"name"`
}

type EDICompanyNotificationRecipientResp struct {
	CompanyNotificationRecipientID mssql.UniqueIdentifier `json:"company_notification_recipient_id"`

	CompanyID mssql.UniqueIdentifier `json:"company_id"`
	// "FORECAST", "ORDER", "INVOICE"
	NotificationType string `json:"notification_type"`

	Email string `json:"email"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type EDICompanyReq struct {
	CompanyID mssql.UniqueIdentifier `json:"company_id"`
	Name      string                 `json:"name"`

	NotificationRecipients []EDICompanyNotificationRecipientReq `json:"notification_recipients"`
}

type EDICompanyNotificationRecipientReq struct {
	CompanyNotificationRecipientID mssql.UniqueIdentifier `json:"company_notification_recipient_id"`

	CompanyID mssql.UniqueIdentifier `json:"company_id"`
	// "FORECAST", "ORDER", "INVOICE"
	NotificationType string `json:"notification_type"`

	Email string `json:"email"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

// ========================================================================================================

type EDIVendorNotificationRecipientResp struct {
	VendorNotificationRecipientID mssql.UniqueIdentifier `json:"vendor_notification_recipient_id"`
	Company                       string                 `json:"company"`
	// "FORECAST", "ORDER", "INVOICE", "VENDOR"
	NotificationType string `json:"notification_type"`

	EDI_PrincipalID mssql.UniqueIdentifier `json:"edi_principal_id"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type EDIVendorNotificationRecipientReq struct {
	VendorNotificationRecipientID mssql.UniqueIdentifier `json:"vendor_notification_recipient_id"`
	Company                       string                 `json:"company"`
	// "FORECAST", "ORDER", "INVOICE", "VENDOR"
	NotificationType string `json:"notification_type"`

	Principal *EDIPrincipalUserEmailReq `json:"principal"`
}

type EDIPrincipalUserEmailReq struct {
	ExternalID string `json:"external_id"`
	Email      string `json:"email"`
}

type EDIVendorNotificationRecipientGroupedResp struct {
	Company   string                    `json:"company"`
	Principal *EDIPrincipalUserEmailReq `json:"principal"`

	// รูปแบบ 1: เก็บเป็น array ของ type
	NotificationTypes []string `json:"notification_types"`
}
