package database

import (
	"fmt"
	"os"

	"github.com/go-ldap/ldap/v3"
)

func AuthenticateUserDomainLogin(user, pass string) (bool, string) {
	ldapServer := os.Getenv("LDAP_IP")
	ldapDNS := os.Getenv("LDAP_DNS")

	if ldapServer == "" {
		return false, "LDAP_IP not defined in env"
	}
	if ldapDNS == "" {
		return false, "LDAP_DNS not defined in env"
	}

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, 389))
	if err != nil {
		return false, fmt.Sprintf("failed to connect to LDAP server: %v", err)
	}
	defer l.Close()

	bindUsername := fmt.Sprintf("%s@%s", user, ldapDNS)

	if err := l.Bind(bindUsername, pass); err != nil {
		return false, fmt.Sprintf("failed to authenticate user %s: %v", user, err)
	}

	return true, ""
}
