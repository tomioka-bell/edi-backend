package ports

import (
	"backend/internal/clients"
	"backend/internal/core/models"
)

type EmployeeService interface {
	GetEmployeeByADLogonService(adUserLogon string) (models.EmployeeView, error)
	StartLoginWithEmailEmployeeOTPService(
		ldapUser *clients.LDAPUserInfo,
	) (string, error)
	CompleteLoginEmployee(account string, code string) (string, error)
}
