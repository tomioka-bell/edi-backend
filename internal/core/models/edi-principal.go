package models

import mssql "github.com/microsoft/go-mssqldb"

type EDIPrincipalUserReq struct {
	EDI_PrincipalID mssql.UniqueIdentifier `json:"edi_principal_id"`
	ExternalID      string                 `json:"external_id"`
	Email           string                 `json:"email"`
	Display_name    string                 `json:"display_name"`
	Profile         string                 `json:"profile"`
	Group           string                 `json:"group"`
	CompanyCode     string                 `json:"company_code"`
	Role            string                 `json:"role"`
	SourceSystem    string                 `json:"source_system"`
	Status          string                 `json:"status"`
	Username        string                 `json:"username"`
	RoleName        string                 `json:"role_name"`
	LoginWithoutOTP bool                   `json:"login_without_otp"`
}

type EDIPrincipalUserByGroup struct {
	EDI_PrincipalID mssql.UniqueIdentifier `json:"edi_principal_id"`
	ExternalID      string                 `json:"external_id"`
	Email           string                 `json:"email"`
	Group           string                 `json:"group"`
	CompanyCode     string                 `json:"company_code"`
	Display_name    string                 `json:"display_name"`
	Username        string                 `json:"username"`
	Role            string                 `json:"role"`
	Status          string                 `json:"status"`
	LoginWithoutOTP bool                   `json:"login_without_otp"`
}

type EDIPrincipalUserByCompany struct {
	EDI_PrincipalID mssql.UniqueIdentifier `json:"edi_principal_id"`
	Display_name    string                 `json:"display_name"`
	ExternalID      string                 `json:"external_id"`
	Email           string                 `json:"email"`
	CompanyCode     string                 `json:"company_code"`
	Status          string                 `json:"status"`
	Group           string                 `json:"group"`
}

type EDIPrincipalUserUpdate struct {
	EDI_PrincipalID string `json:"edi_principal_id"`
	ExternalID      string `json:"external_id"`
	Email           string `json:"email"`
	CompanyCode     string `json:"company_code"`
	Display_name    string `json:"display_name"`
	Profile         string `json:"profile"`
	Group           string `json:"group"`
	Role            string `json:"role"`
	SourceSystem    string `json:"source_system"`
	Status          string `json:"status"`
	Username        string `json:"username"`
	RoleName        string `json:"role_name"`
	LoginWithoutOTP bool   `json:"login_without_otp"`
}
