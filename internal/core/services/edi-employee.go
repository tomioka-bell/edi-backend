package services

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"backend/internal/clients"
	"backend/internal/core/domains"
	"backend/internal/core/models"
	ports "backend/internal/core/ports/repositories"
	servicesports "backend/internal/core/ports/services"
	"backend/internal/pkgs/mailer"
	"backend/internal/pkgs/utils"
)

type employeeService struct {
	employeeRepo ports.EmployeeRepository
}

func NewEmployeeService(EmployeeRepo ports.EmployeeRepository) servicesports.EmployeeService {
	return &employeeService{employeeRepo: EmployeeRepo}
}

func (s *employeeService) GetEmployeeByADLogonService(adUserLogon string) (models.EmployeeView, error) {
	employee, err := s.employeeRepo.GetEmployeeByADLogon(adUserLogon)
	if err != nil {
		return models.EmployeeView{}, err
	}

	employeeView := models.EmployeeView{
		UHR_EmpCode:      employee.UHR_EmpCode,
		UHR_FullNameTh:   employee.UHR_FullNameTh,
		UHR_FullNameEn:   employee.UHR_FullNameEn,
		UHR_Department:   employee.UHR_Department,
		UHR_Position:     employee.UHR_Position,
		UHR_StatusToUse:  employee.UHR_StatusToUse,
		AD_UserLogon:     employee.AD_UserLogon,
		AD_Mail:          employee.AD_Mail,
		AD_Phone:         employee.AD_Phone,
		AD_AccountStatus: employee.AD_AccountStatus,
		UHR_OrgGroup:     employee.UHR_OrgGroup,
		UHR_OrgName:      employee.UHR_OrgName,
	}

	return employeeView, nil
}

func (s *employeeService) StartLoginWithEmailEmployeeOTPService(
	ldapUser *clients.LDAPUserInfo,
) (string, error) {

	if ldapUser == nil {
		return "", errors.New("ldap user is nil")
	}

	// 1) แปลงข้อมูลจาก LDAP → principal
	principal, err := s.ensurePrincipalFromLDAPUser(ldapUser)
	if err != nil {
		return "", fmt.Errorf("ensurePrincipalFromLDAPUser error: %w", err)
	}

	// 2) ถ้า policy อนุญาต login โดยไม่ต้อง OTP
	if principal.LoginWithoutOTP {
		token, err := utils.GenerateJWTFromPrincipal(principal)
		if err != nil {
			return "", fmt.Errorf("GenerateJWTFromPrincipal error: %w", err)
		}
		return token, nil
	}

	// 3) ถ้ายังต้องใช้ OTP → สร้าง/บันทึก OTP และส่งอีเมล
	u, err := s.employeeRepo.StartLoginWithEmailEmployeeOTP(ldapUser.ADMail, ldapUser.ADUsername, ldapUser.EmployeeCode)
	if err != nil {
		return "", fmt.Errorf("StartLoginWithEmailEmployeeOTP error: %w", err)
	}
	if u == nil {
		return "", errors.New("เริ่มกระบวนการ OTP ไม่สำเร็จ (ผลลัพธ์ว่าง)")
	}

	if err := mailer.SendLoginOTPEmail(ldapUser.ADMail, u.TempOTP); err != nil {
		return "", fmt.Errorf("send login otp email: %w", err)
	}

	// ใช้ OTP → ยังไม่สร้าง token
	return "", nil
}

func (s *employeeService) ensurePrincipalFromLDAPUser(emp *clients.LDAPUserInfo) (*domains.EDI_Principal, error) {
	displayName := strings.TrimSpace(emp.FullnameEN)

	principal, err := s.employeeRepo.FindPrincipalByEmail(emp.ADMail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			principal = &domains.EDI_Principal{
				ExternalID:      emp.EmployeeCode,
				SourceSystem:    "APP_EMPLOYEE",
				Email:           emp.ADMail,
				Username:        emp.ADUsername,
				DisplayName:     displayName,
				Profile:         "",
				Role:            "SU",
				Group:           "Prospira (Thailand) Co., Ltd.",
				Status:          "active",
				LoginWithoutOTP: false,
			}

			if err2 := s.employeeRepo.CreateEDIPrincipalRepository(principal); err2 != nil {
				if isUniqueViolation(err2) {
					principal.Username = fmt.Sprintf("%s_%d", emp.ADUsername, time.Now().Unix())
					if err3 := s.employeeRepo.CreateEDIPrincipalRepository(principal); err3 != nil {
						return nil, fmt.Errorf("create principal (retry): %w", err3)
					}
				} else {
					return nil, fmt.Errorf("create principal: %w", err2)
				}
			}
		} else {
			return nil, err
		}
	} else {
		updates := map[string]interface{}{}

		if principal.ExternalID != emp.EmployeeCode {
			updates["external_id"] = emp.EmployeeCode
			principal.ExternalID = emp.EmployeeCode
		}
		if principal.Email != emp.ADMail {
			updates["email"] = emp.ADMail
			principal.Email = emp.ADMail
		}
		if principal.Username != emp.ADUsername {
			updates["username"] = emp.ADUsername
			principal.Username = emp.ADUsername
		}
		if principal.DisplayName != displayName {
			updates["display_name"] = displayName
			principal.DisplayName = displayName
		}
		if principal.Group != "Prospira (Thailand) Co., Ltd." {
			updates["group"] = "Prospira (Thailand) Co., Ltd."
			principal.Group = "Prospira (Thailand) Co., Ltd."
		}
		if principal.Status != "active" {
			updates["status"] = "active"
			principal.Status = "active"
		}
		// if principal.LoginWithoutOTP != false {
		// 	updates["login_without_otp"] = false
		// 	principal.LoginWithoutOTP = false
		// }

		if len(updates) > 0 {
			if err := s.employeeRepo.UpdatePrincipalWithMap(principal.EDI_PrincipalID.String(), updates); err != nil {
				return nil, fmt.Errorf("update principal: %w", err)
			}
		}
	}

	if principal == nil {
		return nil, errors.New("principal is nil after ensure")
	}
	return principal, nil
}

func (s *employeeService) CompleteLoginEmployee(email string, code string) (string, error) {
	// 1) verify code (OTP) – ไม่ยุ่งกับ EmployeeView / LDAP ตรงนี้แล้ว
	_, err := s.employeeRepo.VerifyLoginCodeEmployee(email, code)
	if err != nil {
		return "", err
	}

	// 2) ดึง principal จาก email (สร้างไว้แล้วตอน LDAP login + send OTP)
	principal, err := s.employeeRepo.FindPrincipalByUsername(email)
	if err != nil {
		return "", fmt.Errorf("find principal by email: %w", err)
	}
	if principal == nil {
		return "", errors.New("ไม่พบ principal สำหรับอีเมลนี้ กรุณาลอง login ใหม่")
	}

	// 3) ออก JWT จาก principal
	secret := os.Getenv("TOKEN_SECRET_KEY")
	if secret == "" {
		return "", errors.New("missing TOKEN_SECRET_KEY")
	}
	jwtSecretKey := []byte(secret)

	claims := jwt.MapClaims{
		"user_id":      principal.ExternalID,
		"username":     principal.Username,
		"display_name": principal.DisplayName,
		"profile":      principal.Profile,
		"group":        principal.Group,
		"status":       principal.Status,
		"iat":          time.Now().Unix(),
		"exp":          time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", errors.New("เกิดข้อผิดพลาดในการเซ็นชื่อ JWT")
	}
	return signedToken, nil
}

//=============================================================================================

// func (s *employeeService) CompleteLoginEmployee(email string, code string) (string, error) {

// 	// 1) verify OTP ด้วย email + code
// 	user, err := s.employeeRepo.VerifyLoginCodeEmployee(email, code)
// 	if err != nil {
// 		return "", err
// 	}
// 	if user == nil {
// 		return "", errors.New("verify code สำเร็จ แต่ไม่พบข้อมูลผู้ใช้")
// 	}

// 	// ใช้ user จาก Verify เป็น emp ต่อได้เลย (ไม่ต้อง query ซ้ำ)
// 	emp := user
// 	displayName := strings.TrimSpace(emp.UHR_FullNameEn)

// 	// 2) ลองหา principal จาก email
// 	principal, err := s.employeeRepo.FindPrincipalByEmail(emp.AD_Mail)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			// --------- กรณีไม่พบ principal: สร้างใหม่ ---------
// 			principal = &domains.EDI_Principal{
// 				ExternalID:   emp.UHR_EmpCode,
// 				SourceSystem: "APP_EMPLOYEE",
// 				Email:        emp.AD_Mail,
// 				Username:     emp.AD_UserLogon,
// 				DisplayName:  displayName,
// 				Profile:      "",
// 				Role:         emp.UHR_Position,
// 				Group:        "Prospira (Thailand) Co., Ltd.",
// 				Status:       "active",
// 			}

// 			if err2 := s.employeeRepo.CreateEDIPrincipalRepository(principal); err2 != nil {
// 				if isUniqueViolation(err2) {
// 					// ถ้า username ซ้ำ ให้เปลี่ยนชื่อแล้วลองใหม่
// 					principal.Username = fmt.Sprintf("%s_%d", emp.AD_UserLogon, time.Now().Unix())
// 					if err3 := s.employeeRepo.CreateEDIPrincipalRepository(principal); err3 != nil {
// 						return "", fmt.Errorf("create principal (retry): %w", err3)
// 					}
// 				} else {
// 					return "", fmt.Errorf("create principal: %w", err2)
// 				}
// 			}
// 		} else {
// 			// error อื่น ๆ จาก FindPrincipalByEmail
// 			return "", err
// 		}
// 	} else {
// 		// --------- กรณีพบ principal: ตรวจว่าข้อมูลเปลี่ยนไหม ถ้าเปลี่ยนให้ update ---------
// 		updates := map[string]interface{}{}

// 		if principal.ExternalID != emp.UHR_EmpCode {
// 			updates["external_id"] = emp.UHR_EmpCode
// 			principal.ExternalID = emp.UHR_EmpCode
// 		}
// 		if principal.Email != emp.AD_Mail {
// 			updates["email"] = emp.AD_Mail
// 			principal.Email = emp.AD_Mail
// 		}
// 		if principal.Username != emp.AD_UserLogon {
// 			updates["username"] = emp.AD_UserLogon
// 			principal.Username = emp.AD_UserLogon
// 		}
// 		if principal.DisplayName != displayName {
// 			updates["display_name"] = displayName
// 			principal.DisplayName = displayName
// 		}
// 		if principal.Role != emp.UHR_Position {
// 			updates["role"] = emp.UHR_Position
// 			principal.Role = emp.UHR_Position
// 		}

// 		if principal.Group != "Prospira (Thailand) Co., Ltd." {
// 			updates["group"] = "Prospira (Thailand) Co., Ltd."
// 			principal.Group = "Prospira (Thailand) Co., Ltd."
// 		}
// 		if principal.Status != "active" {
// 			updates["status"] = "active"
// 			principal.Status = "active"
// 		}

// 		if len(updates) > 0 {
// 			if err := s.employeeRepo.UpdatePrincipalWithMap(principal.EDI_PrincipalID.String(), updates); err != nil {
// 				return "", fmt.Errorf("update principal: %w", err)
// 			}
// 		}
// 	}

// 	if principal == nil {
// 		return "", errors.New("principal is nil after ensure")
// 	}

// 	// 3) JWT ใช้ค่าจาก principal (ซึ่งอาจถูก update แล้ว)
// 	secret := os.Getenv("TOKEN_SECRET_KEY")
// 	if secret == "" {
// 		return "", errors.New("missing TOKEN_SECRET_KEY")
// 	}
// 	jwtSecretKey := []byte(secret)

// 	claims := jwt.MapClaims{
// 		"user_id":      principal.ExternalID,
// 		"username":     principal.Username,
// 		"display_name": principal.DisplayName,
// 		"profile":      principal.Profile,
// 		"group":        principal.Group,
// 		"status":       principal.Status,
// 		"iat":          time.Now().Unix(),
// 		"exp":          time.Now().Add(24 * time.Hour).Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	signedToken, err := token.SignedString(jwtSecretKey)
// 	if err != nil {
// 		return "", errors.New("เกิดข้อผิดพลาดในการเซ็นชื่อ JWT")
// 	}
// 	return signedToken, nil
// }
