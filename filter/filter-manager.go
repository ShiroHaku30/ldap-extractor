package filter

import (
	"strings"

	"github.com/go-ldap/ldap/v3"
)

type Manager struct {
	SAMAccountName string
	EmployeeID     string
}

type ManagerLookup map[string]Manager

func BuildManagerLookup(
	entries []*ldap.Entry,
) ManagerLookup {

	lookup := make(ManagerLookup)

	for _, entry := range entries {

		samAccountName := entry.GetAttributeValue(
			"sAMAccountName",
		)

		if samAccountName == "" {
			continue
		}

		lookup[normalizeDN(entry.DN)] = Manager{
			SAMAccountName: samAccountName,
			EmployeeID: entry.GetAttributeValue(
				"employeeID",
			),
		}
	}

	return lookup
}

func FindManager(
	entry *ldap.Entry,
	lookup ManagerLookup,
) (Manager, bool) {

	managerDN := entry.GetAttributeValue(
		"manager",
	)

	if managerDN == "" {
		return Manager{}, false
	}

	manager, ok := lookup[normalizeDN(managerDN)]

	return manager, ok
}

func normalizeDN(dn string) string {
	return strings.ToLower(
		strings.TrimSpace(dn),
	)
}
