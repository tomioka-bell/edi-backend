package repositories

import (
	"backend/internal/core/domains"
	"fmt"

	"gorm.io/gorm"
)

func (r *UserRepositoryDB) CreateEDIPrincipalRepository(EDIUser *domains.EDI_Principal) error {
	if err := r.db.Create(EDIUser).Error; err != nil {
		fmt.Printf("CreateEDIPrincipalRepository error: %v\n", err)
		return err
	}
	return nil
}

func (r *UserRepositoryDB) FindByExternalID(externalID string) (*domains.EDI_Principal, error) {
	var user domains.EDI_Principal

	if err := r.db.Where("external_id = ?", externalID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryDB) FindPrincipalByEmail(email string) (*domains.EDI_Principal, error) {
	var user domains.EDI_Principal

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryDB) FindPrincipalByGroup(group string) ([]domains.EDI_Principal, error) {
	var users []domains.EDI_Principal

	if err := r.db.Where("[group] = ?", group).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoryDB) UpdatePrincipalWithMap(principalID string, updates map[string]interface{}) error {
	return r.db.Model(&domains.EDI_Principal{}).
		Where("edi_principal_id = ?", principalID).
		Updates(updates).
		Error
}
