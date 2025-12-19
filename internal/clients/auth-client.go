package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type LDAPUserInfo struct {
	EmployeeCode string `json:"employee_code"`
	PrefixTH     string `json:"prefix_th"`
	FirstnameTH  string `json:"firstname_th"`
	LastnameTH   string `json:"lastname_th"`
	FullnameTH   string `json:"fullname_th"`
	PrefixEN     string `json:"prefix_en"`
	FirstnameEN  string `json:"firstname_en"`
	LastnameEN   string `json:"lastname_en"`
	FullnameEN   string `json:"fullname_en"`
	Sex          string `json:"sex"`
	Department   string `json:"department"`
	Position     string `json:"position"`
	ADUsername   string `json:"ad_username"`
	ADMail       string `json:"ad_mail"`
	TempOTP      string `gorm:"-"`
}

type ldapAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ldapAuthResponse struct {
	Err      bool           `json:"err"`
	Message  string         `json:"message"`
	UserInfo []LDAPUserInfo `json:"user_info"`
}

// คืนค่า: user (ข้อมูลจาก LDAP ถ้า success), ok, message
func LdapAuthenticate(username, password string) (*LDAPUserInfo, bool, string) {
	baseURL := os.Getenv("AUTH_SERVICE_URL")
	url := fmt.Sprintf("%s/auth/ldap-authen", baseURL)

	reqBody := ldapAuthRequest{
		Username: username,
		Password: password,
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, false, "failed to marshal request to auth-service: " + err.Error()
	}

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, false, "failed to create request to auth-service: " + err.Error()
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, false, "failed to call auth-service: " + err.Error()
	}
	defer resp.Body.Close()

	var res ldapAuthResponse
	_ = json.NewDecoder(resp.Body).Decode(&res)

	msg := res.Message

	if resp.StatusCode != http.StatusOK {
		if msg != "" {
			return nil, false, msg
		}
		return nil, false, fmt.Sprintf("auth-service returned status %d", resp.StatusCode)
	}

	// กรณี service บอก err = true
	if res.Err {
		if msg == "" {
			msg = "authentication failed"
		}
		return nil, false, msg
	}

	if len(res.UserInfo) == 0 {
		return nil, false, "authentication success but no user_info returned"
	}

	return &res.UserInfo[0], true, msg
}
