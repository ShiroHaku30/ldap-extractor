package ldap

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"ldap-extractor/config"
)

func TestConnection(cfg *config.Config) error {
	conn, err := ldap.DialURL(cfg.LDAP.URL)
	if err != nil {
		return fmt.Errorf("LDAP connection failed: %w", err)
	}
	defer conn.Close()

	fmt.Println("LDAP server reachable")

	err = conn.Bind(
		cfg.LDAP.Username,
		cfg.LDAP.Password,
	)

	if err != nil {
		return fmt.Errorf("LDAP bind failed: %w", err)
	}

	fmt.Println("LDAP authentication successful")

	return nil
}
