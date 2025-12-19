package ports

import "backend/internal/core/models"

type EDICompanyService interface {
	CreateCompanyService(req models.EDICompanyResp) error
	CreateNotificationRecipientService(req models.EDICompanyNotificationRecipientResp) error
	GetCompanyByCompanyIDService(companyID string) (*models.EDICompanyReq, error)
	GetEDIVendorNotificationRecipientByCompanyService(company string) ([]models.EDIVendorNotificationRecipientReq, error)
	CreateVendorNotificationRecipientService(req models.EDIVendorNotificationRecipientResp) error
	DeleteNotificationRecipientVendorService(vendorNotificationRecipientID string) error
}
