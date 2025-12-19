package ports

import "backend/internal/core/domains"

type EDICompanyRepository interface {
	CreateCompany(company *domains.EDICompany) error
	CreateNotificationRecipient(recipient *domains.EDICompanyNotificationRecipient) error
	GetCompanyByCompanyID(companyID string) (*domains.EDICompany, error)
	GetNotificationRecipientByCompanyID(companyID string) ([]domains.EDICompanyNotificationRecipient, error)

	CreateNotificationRecipientVendor(Employee *domains.EDIVendorNotificationRecipient) error
	GetEDIVendorNotificationRecipientByCompany(company string) ([]domains.EDIVendorNotificationRecipient, error)
	DeleteNotificationRecipientVendor(vendorNotificationRecipientID string) error
}
