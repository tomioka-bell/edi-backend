package ports

import (
	"backend/internal/core/models"
	"context"
)

type UserService interface {
	CreateUserService(req models.UserResp) error
	SignIn(dto models.LoginResp) (string, error)
	GetProfileByCookieId(userID string) (models.UserReq, error)
	GetAllUserSevice() ([]models.UserReqAll, error)
	UpdateUserWithMapService(userID string, updates map[string]interface{}) error
	GetUserByID(userID string) (models.UserAdminReq, error)
	GetUserCountService() (int64, error)
	RequestPasswordReset(email, lang string) error
	PerformPasswordResetByToken(token, newPassword string) error
	StartLogin(email, password string) error
	CompleteLogin(email, code string) (string, error)
	GetEDIUserByGroupService(group string) ([]models.UserReqAll, error)
	SendInitialPasswordSetup(email, lang string) error
	GetActiveEmailsByGroup(ctx context.Context, group string) ([]string, error)
	RequestPasswordResetMany(emails []string, lang string) error

	// ======================== EDI Principal ==================================================
	GetEDIPrincipalUserByID(ExternalID string) (models.EDIPrincipalUserReq, error)
	GetEDIPrincipalUserByGroupService(group string) ([]models.EDIPrincipalUserByGroup, error)
	GetEDIPrincipalUserByCompanyService(company string) ([]models.EDIPrincipalUserByCompany, error)
	StartLoginWithEmailEmployeeOTPService(email, password string) (string, error)
	UpdatePrincipalWithMapService(ExternalID string, updates map[string]interface{}) error
}
