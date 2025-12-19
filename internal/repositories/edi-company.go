package repositories

import (
	"backend/internal/core/domains"
	ports "backend/internal/core/ports/repositories"
	"fmt"

	"gorm.io/gorm"
)

type EDICompanyRepositoryDB struct {
	db *gorm.DB
}

func NewEDICompanyRepositoryDB(db *gorm.DB) ports.EDICompanyRepository {
	// if err := db.AutoMigrate(&domains.EDIVendorNotificationRecipient{}); err != nil {
	// 	fmt.Printf("failed to auto migrate: %v", err)
	// }
	return &EDICompanyRepositoryDB{db: db}
}

func (r *EDICompanyRepositoryDB) CreateCompany(Employee *domains.EDICompany) error {
	if err := r.db.Create(Employee).Error; err != nil {
		fmt.Printf("CreateCompany error: %v\n", err)
		return err
	}
	return nil
}

func (r *EDICompanyRepositoryDB) CreateNotificationRecipient(Employee *domains.EDICompanyNotificationRecipient) error {
	if err := r.db.Create(Employee).Error; err != nil {
		fmt.Printf("CreateNotificationRecipient error: %v\n", err)
		return err
	}
	return nil
}

func (r *EDICompanyRepositoryDB) GetCompanyByCompanyID(companyID string) (*domains.EDICompany, error) {
	var company domains.EDICompany
	if err := r.db.Where("company_id = ?", companyID).Preload("NotificationRecipients").First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *EDICompanyRepositoryDB) GetNotificationRecipientByCompanyID(companyID string) ([]domains.EDICompanyNotificationRecipient, error) {
	var recipients []domains.EDICompanyNotificationRecipient
	if err := r.db.Where("company_id = ?", companyID).Find(&recipients).Error; err != nil {
		return nil, err
	}
	return recipients, nil
}

func (r *EDICompanyRepositoryDB) CreateNotificationRecipientVendor(Employee *domains.EDIVendorNotificationRecipient) error {
	if err := r.db.Create(Employee).Error; err != nil {
		fmt.Printf("CreateNotificationRecipientVendor error: %v\n", err)
		return err
	}
	return nil
}

func (r *EDICompanyRepositoryDB) GetEDIVendorNotificationRecipientByCompany(company string) ([]domains.EDIVendorNotificationRecipient, error) {
	var recipients []domains.EDIVendorNotificationRecipient

	err := r.db.
		Model(&domains.EDIVendorNotificationRecipient{}).
		Where("company = ? AND deleted_at IS NULL", company).
		Preload("Principal", "deleted_at IS NULL").
		Find(&recipients).Error

	if err != nil {
		return nil, err
	}
	return recipients, nil
}

// func (r *EDICompanyRepositoryDB) GetEDIVendorNotificationRecipientByCompany(company string) ([]domains.EDIVendorNotificationRecipient, error) {
// 	fmt.Println("company : ", company)

// 	rows, err := r.db.Debug().Raw(`
//     SELECT
//         r.vendor_notification_recipient_id,
//         r.company,
//         r.notification_type,
//         r.edi_principal_id,
//         r.created_at,
//         r.updated_at,
//         r.deleted_at,
// 		  ISNULL(p.edi_principal_id, CAST('00000000-0000-0000-0000-000000000000' AS uniqueidentifier)) AS p_edi_principal_id,

//         p.external_id AS p_external_id,
//         p.source_system AS p_source_system,
//         p.email AS p_email,
//         p.display_name AS p_display_name,
//         p.profile AS p_profile,
//         p.[group] AS p_group,
//         p.role AS p_role,
//         p.status AS p_status,
//         p.username AS p_username,
//         p.created_at AS p_created_at,
//         p.updated_at AS p_updated_at,
//         p.deleted_at AS p_deleted_at
//     FROM edi_vendor_notification_recipient r
//     LEFT JOIN edi_principals p
//         ON r.edi_principal_id = p.edi_principal_id
//     WHERE r.company = ?
//       AND r.deleted_at IS NULL
// `, company).Rows()

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	recipients := []domains.EDIVendorNotificationRecipient{}

// 	for rows.Next() {
// 		var rec domains.EDIVendorNotificationRecipient
// 		var principal domains.EDI_Principal

// 		err := rows.Scan(
// 			&rec.VendorNotificationRecipientID,
// 			&rec.Company,
// 			&rec.NotificationType,
// 			&rec.EDI_PrincipalID,
// 			&rec.CreatedAt,
// 			&rec.UpdatedAt,
// 			&rec.DeletedAt,

// 			&principal.EDI_PrincipalID,
// 			&principal.ExternalID,
// 			&principal.SourceSystem,
// 			&principal.Email,
// 			&principal.DisplayName,
// 			&principal.Profile,
// 			&principal.Group,
// 			&principal.Role,
// 			&principal.Status,
// 			&principal.Username,
// 			&principal.CreatedAt,
// 			&principal.UpdatedAt,
// 			&principal.DeletedAt,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 		if principal.EDI_PrincipalID != (mssql.UniqueIdentifier{}) {
// 			rec.Principal = &principal
// 		}

// 		recipients = append(recipients, rec)
// 	}

// 	return recipients, nil
// }

func (r *EDICompanyRepositoryDB) DeleteNotificationRecipientVendor(vendorNotificationRecipientID string) error {
	if err := r.db.
		Unscoped().
		Where("vendor_notification_recipient_id = ?", vendorNotificationRecipientID).
		Delete(&domains.EDIVendorNotificationRecipient{}).Error; err != nil {

		fmt.Printf("DeleteNotificationRecipientVendor error: %v\n", err)
		return err
	}
	return nil
}
