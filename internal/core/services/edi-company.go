package services

import (
	"backend/internal/core/domains"
	"backend/internal/core/models"
	ports "backend/internal/core/ports/repositories"
	servicesports "backend/internal/core/ports/services"
	"fmt"

	"github.com/google/uuid"
	mssql "github.com/microsoft/go-mssqldb"
)

type EDICompanyService struct {
	companyRepo ports.EDICompanyRepository
}

func NewEDICompanyService(companyRepo ports.EDICompanyRepository) servicesports.EDICompanyService {
	return &EDICompanyService{companyRepo: companyRepo}
}

func (s *EDICompanyService) CreateCompanyService(req models.EDICompanyResp) error {
	u := uuid.New()

	var newID mssql.UniqueIdentifier
	copy(newID[:], u[:])

	domainISR := domains.EDICompany{
		CompanyID: newID,
		Name:      req.Name,
	}

	if s.companyRepo == nil {
		return fmt.Errorf("company repository is not initialized")
	}
	if err := s.companyRepo.CreateCompany(&domainISR); err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}
	return nil
}

func (s *EDICompanyService) CreateNotificationRecipientService(resp models.EDICompanyNotificationRecipientResp) error {
	u := uuid.New()

	var newID mssql.UniqueIdentifier
	copy(newID[:], u[:])

	domainISR := domains.EDICompanyNotificationRecipient{
		CompanyNotificationRecipientID: newID,
		CompanyID:                      resp.CompanyID,
		NotificationType:               resp.NotificationType,
		Email:                          resp.Email,
	}

	if s.companyRepo == nil {
		return fmt.Errorf("company repository is not initialized")
	}
	if err := s.companyRepo.CreateNotificationRecipient(&domainISR); err != nil {
		return fmt.Errorf("failed to create company notification recipient: %w", err)
	}
	return nil
}

func (s *EDICompanyService) GetCompanyByCompanyIDService(companyID string) (*models.EDICompanyReq, error) {
	company, err := s.companyRepo.GetCompanyByCompanyID(companyID)
	if err != nil {
		return nil, err
	}

	companyReq := models.EDICompanyReq{
		CompanyID: company.CompanyID,
		Name:      company.Name,
	}

	return &companyReq, nil
}

func (s *EDICompanyService) CreateVendorNotificationRecipientService(req models.EDIVendorNotificationRecipientResp) error {
	u := uuid.New()

	var newID mssql.UniqueIdentifier
	copy(newID[:], u[:])

	domainISR := domains.EDIVendorNotificationRecipient{
		VendorNotificationRecipientID: newID,
		Company:                       req.Company,
		NotificationType:              req.NotificationType,
		EDI_PrincipalID:               req.EDI_PrincipalID,
	}

	if s.companyRepo == nil {
		return fmt.Errorf("company is not initialized")
	}
	if err := s.companyRepo.CreateNotificationRecipientVendor(&domainISR); err != nil {
		return fmt.Errorf("failed to create vendor notification recipient: %w", err)
	}
	return nil
}

func (s *EDICompanyService) GetEDIVendorNotificationRecipientByCompanyService(company string) ([]models.EDIVendorNotificationRecipientReq, error) {
	recipients, err := s.companyRepo.GetEDIVendorNotificationRecipientByCompany(company)
	if err != nil {
		return nil, err
	}

	res := make([]models.EDIVendorNotificationRecipientReq, 0, len(recipients))

	for _, r := range recipients {
		var principal *models.EDIPrincipalUserEmailReq
		if r.Principal != nil {
			principal = &models.EDIPrincipalUserEmailReq{
				ExternalID: r.Principal.ExternalID,
				Email:      r.Principal.Email,
			}
		}

		res = append(res, models.EDIVendorNotificationRecipientReq{
			VendorNotificationRecipientID: r.VendorNotificationRecipientID,
			Company:                       r.Company,
			NotificationType:              r.NotificationType,
			Principal:                     principal,
		})
	}

	return res, nil
}

func (s *EDICompanyService) DeleteNotificationRecipientVendorService(vendorNotificationRecipientID string) error {
	if s.companyRepo == nil {
		return fmt.Errorf("company is not initialized")
	}
	if err := s.companyRepo.DeleteNotificationRecipientVendor(vendorNotificationRecipientID); err != nil {
		return fmt.Errorf("failed to delete vendor notification recipient: %w", err)
	}
	return nil
}
