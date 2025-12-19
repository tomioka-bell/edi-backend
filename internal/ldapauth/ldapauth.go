package ldapauth

import (
	"fmt"
	"os"

	"github.com/go-ldap/ldap/v3"
)

type Service interface {
	Authenticate(user, pass string) (bool, string)
}

type service struct {
	server string
	dns    string
}

func NewService() (Service, error) {
	ldapServer := os.Getenv("LDAP_IP")
	ldapDNS := os.Getenv("LDAP_DNS")

	if ldapServer == "" {
		return nil, fmt.Errorf("LDAP_IP not defined in env")
	}
	if ldapDNS == "" {
		return nil, fmt.Errorf("LDAP_DNS not defined in env")
	}

	return &service{
		server: ldapServer,
		dns:    ldapDNS,
	}, nil
}

func (s *service) Authenticate(user, pass string) (bool, string) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.server, 389))
	if err != nil {
		return false, fmt.Sprintf("failed to connect to LDAP server: %v", err)
	}
	defer l.Close()

	bindUsername := fmt.Sprintf("%s@%s", user, s.dns)

	if err := l.Bind(bindUsername, pass); err != nil {
		return false, fmt.Sprintf("failed to authenticate user %s: %v", user, err)
	}

	return true, ""
}
