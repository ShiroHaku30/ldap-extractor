package ldap

import (
	"fmt"
	"os"

	"github.com/go-ldap/ldap/v3"

	"ldap-extractor/config"
)

func Dump(cfg *config.Config) error {
	conn, err := ldap.DialURL(cfg.LDAP.URL)
	if err != nil {
		return fmt.Errorf("LDAP connection failed: %w", err)
	}
	defer conn.Close()

	// Authenticate
	if err := conn.Bind(
		cfg.LDAP.Username,
		cfg.LDAP.Password,
	); err != nil {
		return fmt.Errorf("LDAP bind failed: %w", err)
	}

	fmt.Println("LDAP bind successful")

	// Create output file
	file, err := os.Create(cfg.Output.File)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	pageSize := cfg.LDAP.Search.PageSize

	pagingControl := ldap.NewControlPaging(pageSize)

	total := 0
	page := 0

	for {
		page++

		searchRequest := ldap.NewSearchRequest(
			cfg.LDAP.BaseDN,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0, // Size limit
			0, // Time limit
			false,
			cfg.LDAP.Search.Filter,
			cfg.LDAP.Search.Attributes,
			[]ldap.Control{pagingControl},
		)

		result, err := conn.Search(searchRequest)
		if err != nil {
			return fmt.Errorf(
				"LDAP search failed on page %d: %w",
				page,
				err,
			)
		}

		fmt.Printf(
			"Page %d: %d entries\n",
			page,
			len(result.Entries),
		)

		for _, entry := range result.Entries {
			if err := writeEntry(file, entry); err != nil {
				return fmt.Errorf(
					"failed writing entry: %w",
					err,
				)
			}

			total++
		}

		// Flush the file after every page.
		if err := file.Sync(); err != nil {
			return fmt.Errorf(
				"failed to flush output file: %w",
				err,
			)
		}

		// Get the paging control returned by the server.
		updatedControl := ldap.FindControl(
			result.Controls,
			ldap.ControlTypePaging,
		)

		if updatedControl == nil {
			break
		}

		ctrl, ok := updatedControl.(*ldap.ControlPaging)
		if !ok {
			return fmt.Errorf(
				"invalid LDAP paging control returned by server",
			)
		}

		// Empty cookie means there are no more results.
		if len(ctrl.Cookie) == 0 {
			break
		}

		// Use the cookie for the next page.
		pagingControl.SetCookie(ctrl.Cookie)
	}

	fmt.Printf(
		"LDAP dump complete: %d entries\n",
		total,
	)

	return nil
}

func writeEntry(file *os.File, entry *ldap.Entry) error {
	if _, err := fmt.Fprintf(
		file,
		"dn: %s\n",
		entry.DN,
	); err != nil {
		return err
	}

	for _, attr := range entry.Attributes {
		for _, value := range attr.Values {
			if _, err := fmt.Fprintf(
				file,
				"%s: %s\n",
				attr.Name,
				value,
			); err != nil {
				return err
			}
		}
	}

	_, err := fmt.Fprintln(file)

	return err
}
