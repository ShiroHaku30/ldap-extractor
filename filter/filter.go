package filter

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/go-ldap/ldif"
)

type User struct {
	SAMAccountName string `json:"sAMAccountName"`
	Department     string `json:"department"`
	EmployeeID     string `json:"employeeID"`
}

type Attribute struct {
	Name     string
	Value    string
	IsBase64 bool
}

func encodeJSONValue(value string, encode bool) string {

	if !encode {
		return value
	}

	return base64.StdEncoding.EncodeToString(
		[]byte(value),
	)
}

func ExportJSON(
	input string,
	output string,
	attributes []string,
) error {

	file, err := os.Open(input)
	if err != nil {
		return err
	}

	defer file.Close()

	var data ldif.LDIF

	err = ldif.Unmarshal(
		file,
		&data,
	)

	if err != nil {
		return err
	}

	users := make([]User, 0)

	for _, entry := range data.AllEntries() {

		user := User{
			SAMAccountName: entry.GetAttributeValue(
				"sAMAccountName",
			),

			Department: encodeJSONValue(
				entry.GetAttributeValue("department"),
				true,
			),

			EmployeeID: entry.GetAttributeValue(
				"employeeID",
			),
		}

		// skip entries without username
		if user.SAMAccountName == "" {
			continue
		}

		users = append(users, user)
	}

	out, err := os.Create(output)
	if err != nil {
		return err
	}

	defer out.Close()

	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")

	return encoder.Encode(users)
}
