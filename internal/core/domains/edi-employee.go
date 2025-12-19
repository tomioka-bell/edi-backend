package domains

import (
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"gorm.io/gorm"
)

// {
//     "uhr_emp_code": "P250170",
//     "uhr_full_name_th": "นายภูวดล อาจมิตร์",
//     "uhr_department": "Information Technology",
//     "uhr_position": "ENGINEER 1",
//     "uhr_status_to_use": "ENABLE",
//     "ad_user_logon": "pva",
//     "ad_mail": "puvadon.artmit@prospira.com",
//     "ad_phone": "3076",
//     "ad_account_status": "ENABLE",
//     "uhr_org_group": "OFFLINE",
//     "uhr_org_name": "SG&A"
// }

type EmployeeUser struct {
	EmployeeID mssql.UniqueIdentifier `gorm:"type:uniqueidentifier;primaryKey;default:(newid())" json:"employee_id"`
	Username   string                 `gorm:"column:username"`
	EmpCode    string                 `gorm:"column:emp_code"`
	FullName   string                 `gorm:"column:full_name"`
	Department string                 `gorm:"column:department"`
	Email      string                 `gorm:"column:email"`
	Position   string                 `gorm:"column:position"`
	Status     string                 `gorm:"column:status"`
	CreatedAt  *time.Time             `json:"created_at"`
	UpdatedAt  *time.Time             `json:"updated_at"`
	DeletedAt  gorm.DeletedAt         `gorm:"index"`
}

type EmployeeView struct {
	UHR_EmpCode      string `gorm:"column:UHR_EmpCode"`
	UHR_FullNameTh   string `gorm:"column:UHR_FullName_th"`
	UHR_FullNameEn   string `gorm:"column:UHR_FullName_en"`
	UHR_Department   string `gorm:"column:UHR_Department"`
	UHR_Position     string `gorm:"column:UHR_Position"`
	UHR_StatusToUse  string `gorm:"column:UHR_StatusToUse"`
	AD_UserLogon     string `gorm:"column:AD_UserLogon"`
	AD_Mail          string `gorm:"column:AD_Mail"`
	AD_Phone         string `gorm:"column:AD_Phone"`
	AD_AccountStatus string `gorm:"column:AD_AccountStatus"`
	UHR_OrgGroup     string `gorm:"column:UHR_OrgGroup"`
	UHR_OrgName      string `gorm:"column:UHR_OrgName"`

	TempOTP string `gorm:"-"`
}

func (EmployeeView) TableName() string { return "V_Employees" }

type LoginVerificationEmployee struct {
	ID           mssql.UniqueIdentifier `gorm:"column:id;type:uniqueidentifier;primaryKey"`
	UserLogin    string                 `gorm:"column:user_login;type:nvarchar(100);index"` // AD_UserLogon
	EmpCode      string                 `gorm:"column:emp_code;type:nvarchar(50)"`
	Email        string                 `gorm:"column:email;type:nvarchar(255)"`
	CodeHash     string                 `gorm:"column:code_hash"`
	ExpiresAt    time.Time              `gorm:"column:expires_at"`
	CreatedAt    time.Time              `gorm:"column:created_at"`
	ConsumedAt   *time.Time             `gorm:"column:consumed_at"`
	AttemptCount int                    `gorm:"column:attempt_count"`
}

func (LoginVerificationEmployee) TableName() string { return "login_verifications_employee" }
