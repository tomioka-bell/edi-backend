package models

import (
	mssql "github.com/microsoft/go-mssqldb"
)

type UserResp struct {
	UserID          string `json:"user_id"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname" `
	Username        string `json:"username"`
	Email           string `json:"email"`
	Profile         string `json:"profile"`
	Group           string `json:"group"`
	Password        string `json:"password"`
	Status          string `json:"status"`
	Role            string `json:"role"`
	CompanyCode     string `json:"company_code"`
	LoginWithoutOTP bool   `json:"login_without_otp"`
}

type RoleInfo struct {
	Name        string `json:"role_name"`
	Description string `json:"role_description"`
}

type UserReq struct {
	UserID      string `json:"user_id"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Profile     string `json:"profile"`
	Group       string `json:"group"`
	Status      string `json:"status"`
	RoleName    string `json:"role_name"`
	CompanyCode string `json:"company_code"`
}

type UserAdminReq struct {
	UserID      mssql.UniqueIdentifier `json:"user_id"`
	Firstname   string                 `json:"firstname"`
	Lastname    string                 `json:"lastname"`
	Username    string                 `json:"username"`
	Email       string                 `json:"email"`
	Password    string                 `json:"password"`
	Status      string                 `json:"status"`
	RoleName    string                 `json:"role_name"`
	CompanyCode string                 `json:"company_code"`
	Group       string                 `json:"group"`
}

type UserReqAll struct {
	UserID          mssql.UniqueIdentifier `json:"user_id"`
	Firstname       string                 `json:"firstname"`
	Lastname        string                 `json:"lastname" `
	Username        string                 `json:"username"`
	Email           string                 `json:"email"`
	Status          string                 `json:"status"`
	Profile         string                 `json:"profile"`
	Group           string                 `json:"group"`
	CompanyCode     string                 `json:"company_code"`
	LoginWithoutOTP bool                   `json:"login_without_otp"`
}

type LoginResp struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginCookieResp struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
