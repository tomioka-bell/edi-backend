package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"backend/internal/core/domains"
	"backend/internal/core/models"
	"backend/internal/pkgs/mailer"
	"backend/internal/pkgs/utils"
)

func (s *userService) GetEDIPrincipalUserByID(ExternalID string) (models.EDIPrincipalUserReq, error) {
	user, err := s.userisrRepo.FindByExternalID(ExternalID)
	if err != nil {
		return models.EDIPrincipalUserReq{}, err
	}

	userReq := models.EDIPrincipalUserReq{
		EDI_PrincipalID: user.EDI_PrincipalID,
		ExternalID:      user.ExternalID,
		Email:           user.Email,
		Display_name:    user.DisplayName,
		Profile:         user.Profile,
		Group:           user.Group,
		CompanyCode:     user.CompanyCode,
		Role:            user.Role,
		SourceSystem:    user.SourceSystem,
		Status:          user.Status,
		Username:        user.Username,
		RoleName:        user.Role,
	}

	return userReq, nil
}

func (s *userService) UpdatePrincipalWithMapService(ExternalID string, updates map[string]interface{}) error {
	return s.userisrRepo.UpdatePrincipalWithMap(ExternalID, updates)
}

func (s *userService) GetEDIPrincipalUserByGroupService(group string) ([]models.EDIPrincipalUserByGroup, error) {
	users, err := s.userisrRepo.FindPrincipalByGroup(group)
	if err != nil {
		return nil, err
	}

	res := make([]models.EDIPrincipalUserByGroup, 0, len(users))

	for _, u := range users {
		res = append(res, models.EDIPrincipalUserByGroup{
			EDI_PrincipalID: u.EDI_PrincipalID,
			ExternalID:      u.ExternalID,
			Email:           u.Email,
			Group:           u.Group,
			CompanyCode:     u.CompanyCode,
			Role:            u.Role,
			Status:          u.Status,
			Username:        u.Username,
			LoginWithoutOTP: u.LoginWithoutOTP,
			Display_name:    u.DisplayName,
		})
	}

	return res, nil
}

func (s *userService) GetEDIPrincipalUserByCompanyService(company string) ([]models.EDIPrincipalUserByCompany, error) {
	users, err := s.userisrRepo.FindPrincipalByGroup(company)
	if err != nil {
		return nil, err
	}

	res := make([]models.EDIPrincipalUserByCompany, 0, len(users))

	for _, u := range users {
		res = append(res, models.EDIPrincipalUserByCompany{
			EDI_PrincipalID: u.EDI_PrincipalID,
			Display_name:    u.DisplayName,
			Status:          u.Status,
			ExternalID:      u.ExternalID,
			Email:           u.Email,
			Group:           u.Group,
			CompanyCode:     u.CompanyCode,
		})
	}

	return res, nil
}

func (s *userService) StartLoginWithEmailEmployeeOTPService(email, password string) (string, error) {
	userData, err := s.userisrRepo.FindByEmail(email)
	if err != nil || userData == nil {
		return "", errors.New("No user found in the system.")
	}
	if userData.Status == "disable" {
		return "", errors.New("This account has been deactivated.")
	}
	if !utils.Compare(password, userData.Password) {
		return "", errors.New("The email or password is incorrect.")
	}

	principal, err := s.ensurePrincipalFromUser(userData)
	if err != nil {
		return "", err
	}

	if principal.LoginWithoutOTP {
		token, err := utils.GenerateJWTFromPrincipal(principal)
		if err != nil {
			return "", err
		}
		return token, nil
	}

	u, err := s.userisrRepo.StartLoginWithEmailOTP(userData.Email)
	if err != nil {
		return "", fmt.Errorf("StartLoginWithEmailEmployeeOTP error: %w", err)
	}
	if u == nil {
		return "", errors.New("เริ่มกระบวนการ OTP ไม่สำเร็จ (ผลลัพธ์ว่าง)")
	}

	fmt.Println("sending OTP to:", u.Email, "code:", u.TempOTP)
	if err := mailer.SendLoginOTPEmail(u.Email, u.TempOTP); err != nil {
		return "", fmt.Errorf("send login otp email: %w", err)
	}

	return "", nil
}

func (s *userService) ensurePrincipalFromUser(user *domains.User) (*domains.EDI_Principal, error) {
	displayName := strings.TrimSpace(user.Firstname + " " + user.Lastname)

	principal, err := s.userisrRepo.FindPrincipalByEmail(user.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// สร้างใหม่
			principal = &domains.EDI_Principal{
				ExternalID:   user.UserID.String(),
				SourceSystem: "APP_EMPLOYEE",
				Email:        user.Email,
				Username:     user.Username,
				DisplayName:  displayName,
				Profile:      "",
				Role:         user.Role,
				Group:        user.Group,
				CompanyCode:  user.CompanyCode,
				Status:       "active",
			}

			if err2 := s.userisrRepo.CreateEDIPrincipalRepository(principal); err2 != nil {
				if isUniqueViolation(err2) {
					principal.Username = fmt.Sprintf("%s_%d", user.Username, time.Now().Unix())
					if err3 := s.userisrRepo.CreateEDIPrincipalRepository(principal); err3 != nil {
						return nil, fmt.Errorf("create principal (retry): %w", err3)
					}
				} else {
					return nil, fmt.Errorf("create principal: %w", err2)
				}
			}
		} else {
			return nil, fmt.Errorf("find principal: %w", err)
		}
	} else {
		updates := map[string]interface{}{}

		if principal.ExternalID != user.UserID.String() {
			updates["external_id"] = user.UserID.String()
			principal.ExternalID = user.UserID.String()
		}
		if principal.Email != user.Email {
			updates["email"] = user.Email
			principal.Email = user.Email
		}
		if principal.Username != user.Username {
			updates["username"] = user.Username
			principal.Username = user.Username
		}
		if principal.DisplayName != displayName {
			updates["display_name"] = displayName
			principal.DisplayName = displayName
		}
		if principal.Role != user.Role {
			updates["role"] = user.Role
			principal.Role = user.Role
		}
		if principal.Group != user.Group {
			updates["group"] = user.Group
			principal.Group = user.Group
		}
		if principal.Status != "active" {
			updates["status"] = "active"
			principal.Status = "active"
		}
		if principal.Profile != user.Profile {
			updates["profile"] = user.Profile
			principal.Profile = user.Profile
		}
		if principal.CompanyCode != user.CompanyCode {
			updates["company_code"] = user.CompanyCode
			principal.CompanyCode = user.CompanyCode
		}
		if principal.LoginWithoutOTP != user.LoginWithoutOTP {
			updates["login_without_otp"] = user.LoginWithoutOTP
			principal.LoginWithoutOTP = user.LoginWithoutOTP
		}

		if len(updates) > 0 {
			if err := s.userisrRepo.UpdatePrincipalWithMap(principal.EDI_PrincipalID.String(), updates); err != nil {
				return nil, fmt.Errorf("update principal: %w", err)
			}
		}
	}

	if principal == nil {
		return nil, errors.New("principal is nil after ensure")
	}
	return principal, nil
}
