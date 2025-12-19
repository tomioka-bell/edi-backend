package ports

import "backend/internal/core/domains"

type UserRepository interface {
	CreateUserRepository(User *domains.User) error
	FindByUsername(username string) (*domains.User, error)
	GetUserByID(userID string) (domains.User, error)
	GetAllUser() ([]domains.User, error)
	UpdateUserWithMap(userID string, updates map[string]interface{}) error
	GetUserCount() (int64, error)
	FindByEmail(email string) (*domains.User, error)
	UpdatePasswordByEmail(email, newPassword string) error
	CreatePasswordResetLink(email, frontendBaseURL, lang string) (string, error)
	ResetPasswordByToken(plainToken, newHashedPassword string) error
	StartLoginWithEmailOTP(email string) (*domains.User, error)
	VerifyLoginCode(email, plainCode string) (*domains.User, error)
	FindUserByGroup(group string) ([]domains.User, error)
	FindActiveEmailsByGroup(group string) ([]string, error)

	// ======================== EDI Principal ==================================================
	FindPrincipalByEmail(email string) (*domains.EDI_Principal, error)
	CreateEDIPrincipalRepository(EDIUser *domains.EDI_Principal) error
	FindByExternalID(externalID string) (*domains.EDI_Principal, error)
	FindPrincipalByGroup(group string) ([]domains.EDI_Principal, error)
	UpdatePrincipalWithMap(principalID string, updates map[string]interface{}) error
}
