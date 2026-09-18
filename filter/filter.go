package filter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-ldap/ldif"
)

type User struct {
	SAMAccountName    string `json:"sAMAccountName"`
	Department        string `json:"department"`
	EmployeeID        string `json:"employeeID"`
	Title             string `json:"title"`
	Manager           string `json:"manager"`
	ManagerEmployeeID string `json:"managerEmployeeID"`
}

func ExportJSON(
	input string,
	output string,
	attributes []string,
) error {

	// Open LDIF file.
	file, err := os.Open(input)
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse LDIF.
	var data ldif.LDIF

	if err := ldif.Unmarshal(file, &data); err != nil {
		return fmt.Errorf(
			"failed to parse LDIF: %w",
			err,
		)
	}

	/*
		Build manager lookup from all LDAP entries.

		The lookup is:

			manager DN
			    ↓
			sAMAccountName
			employeeID

		This is built once and reused for every user.
	*/
	managerLookup := BuildManagerLookup(
		data.AllEntries(),
	)

	users := make([]User, 0)

	/*
		Process each LDAP entry.
	*/
	for _, entry := range data.AllEntries() {

		user := User{
			SAMAccountName: entry.GetAttributeValue(
				"sAMAccountName",
			),

			Department: entry.GetAttributeValue(
				"department",
			),

			EmployeeID: entry.GetAttributeValue(
				"employeeID",
			),

			Title: entry.GetAttributeValue(
				"title",
			),
		}

		// Skip entries without a username.
		if user.SAMAccountName == "" {
			continue
		}

		/*
			Resolve the user's manager.

			The manager attribute has already been
			parsed by ldif.Unmarshal(), so FindManager()
			works with the unmarshalled DN.
		*/
		if manager, ok := FindManager(
			entry,
			managerLookup,
		); ok {

			user.Manager = manager.SAMAccountName
			user.ManagerEmployeeID = manager.EmployeeID
		}

		users = append(users, user)
	}

	// Create JSON output.
	out, err := os.Create(output)
	if err != nil {
		return fmt.Errorf(
			"failed to create JSON output: %w",
			err,
		)
	}
	defer out.Close()

	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(users); err != nil {
		return fmt.Errorf(
			"failed to encode JSON: %w",
			err,
		)
	}

	return nil
}
