package ports

import (
	"backend/internal/clients"
	"backend/internal/core/domains"
)

type EmployeeRepository interface {
	GetEmployeeByADLogon(adUserLogon string) (*domains.EmployeeView, error)
	FindEmployeeByAccount(account string) (*domains.EmployeeView, error)
	FindPrincipalByEmail(email string) (*domains.EDI_Principal, error)
	CreateEDIPrincipalRepository(EDIUser *domains.EDI_Principal) error
	CreateEmployeeRepository(Employee *domains.EmployeeView) error
	StartLoginWithEmailEmployeeOTP(email, ADUsername, EmployeeCode string) (*clients.LDAPUserInfo, error)
	VerifyLoginCodeEmployee(login string, plainCode string) (*domains.LoginVerificationEmployee, error)
	UpdatePrincipalWithMap(principalID string, updates map[string]interface{}) error
	FindPrincipalByUsername(username string) (*domains.EDI_Principal, error)
}
