package models

type EmployeeView struct {
	UHR_EmpCode      string `json:"uhr_emp_code"`
	UHR_FullNameTh   string `json:"uhr_full_name_th"`
	UHR_FullNameEn   string `json:"uhr_full_name_en"`
	UHR_Department   string `json:"uhr_department"`
	UHR_Position     string `json:"uhr_position"`
	UHR_StatusToUse  string `json:"uhr_status_to_use"`
	AD_UserLogon     string `json:"ad_user_logon"`
	AD_Mail          string `json:"ad_mail"`
	AD_Phone         string `json:"ad_phone"`
	AD_AccountStatus string `json:"ad_account_status"`
	UHR_OrgGroup     string `json:"uhr_org_group"`
	UHR_OrgName      string `json:"uhr_org_name"`
}
