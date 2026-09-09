package ldap

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/go-ldap/ldap/v3"

	"ldap-extractor/config"
)

func Dump(cfg *config.Config) error {

	conn, err := ldap.DialURL(cfg.LDAP.URL)
	if err != nil {
		return fmt.Errorf("LDAP connection failed: %w", err)
	}
	defer conn.Close()

	if err := conn.Bind(
		cfg.LDAP.Username,
		cfg.LDAP.Password,
	); err != nil {
		return fmt.Errorf("LDAP bind failed: %w", err)
	}

	fmt.Println("LDAP bind successful")

	dir := filepath.Dir(cfg.Output.File)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf(
			"failed to create output directory: %w",
			err,
		)
	}

	file, err := os.Create(cfg.Output.File)
	if err != nil {
		return fmt.Errorf(
			"failed to create output file: %w",
			err,
		)
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
			0,
			0,
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

			err := writeEntry(
				file,
				entry,
				cfg.LDAP.LDIF.EncodeNonASCII,
			)

			if err != nil {
				return fmt.Errorf(
					"failed writing entry: %w",
					err,
				)
			}

			total++
		}

		if err := file.Sync(); err != nil {
			return fmt.Errorf(
				"failed flushing output: %w",
				err,
			)
		}

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
				"invalid LDAP paging control",
			)
		}

		if len(ctrl.Cookie) == 0 {
			break
		}

		pagingControl.SetCookie(ctrl.Cookie)
	}

	fmt.Printf(
		"LDAP dump complete: %d entries\n",
		total,
	)

	return nil
}

func writeEntry(
	file *os.File,
	entry *ldap.Entry,
	encodeNonASCII bool,
) error {

	writeLDIFValue(
		file,
		"dn",
		[]byte(entry.DN),
		encodeNonASCII,
	)

	for _, attr := range entry.Attributes {

		for _, value := range attr.ByteValues {

			writeLDIFValue(
				file,
				attr.Name,
				value,
				encodeNonASCII,
			)
		}
	}

	fmt.Fprintln(file)

	return nil
}

func writeLDIFValue(
	file *os.File,
	name string,
	value []byte,
	encodeNonASCII bool,
) {

	if needsBase64(
		value,
		encodeNonASCII,
	) {

		fmt.Fprintf(
			file,
			"%s:: %s\n",
			name,
			base64.StdEncoding.EncodeToString(value),
		)

		return
	}

	fmt.Fprintf(
		file,
		"%s: %s\n",
		name,
		string(value),
	)
}

func needsBase64(
	value []byte,
	encodeNonASCII bool,
) bool {

	if len(value) == 0 {
		return false
	}

	// Invalid UTF-8
	if !utf8.Valid(value) {
		return true
	}

	// Optional policy:
	// encode Chinese/Japanese/Korean/etc.
	if encodeNonASCII {

		for _, r := range string(value) {

			if r > 127 {
				return true
			}
		}
	}

	// LDIF SAFE-INIT restrictions
	switch value[0] {

	case ' ', ':', '<':
		return true
	}

	// Unsafe control characters
	for _, b := range value {

		switch b {

		case 0x00: // NULL
			return true

		case 0x0A: // LF
			return true

		case 0x0D: // CR
			return true
		}
	}

	// Recommended by LDIF spec
	// encode trailing spaces
	if value[len(value)-1] == ' ' {
		return true
	}

	return false
}
