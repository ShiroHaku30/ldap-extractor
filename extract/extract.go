package extract

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ldap-extractor/config"
	"ldap-extractor/diff"
	"ldap-extractor/filter"
	"ldap-extractor/ldap"
)

func Run(cfg *config.Config) error {

	timestamp := time.Now().
		Format("20060102_150405")

	fmt.Println(
		"Starting LDAP dump",
	)

	if err := ldap.Dump(cfg); err != nil {
		return fmt.Errorf(
			"dump failed: %w",
			err,
		)
	}

	if err := os.MkdirAll(
		cfg.Output.Filter.Directory,
		0755,
	); err != nil {
		return err
	}

	filteredFile := filepath.Join(
		cfg.Output.Filter.Directory,
		fmt.Sprintf(
			"%s-%s",
			timestamp,
			cfg.Output.Filter.Suffix,
		),
	)

	fmt.Println(
		"Generating filtered output:",
		filteredFile,
	)

	if err := filter.ExportJSON(
		cfg.Output.File,
		filteredFile,
		cfg.LDAP.Search.Attributes,
	); err != nil {
		return fmt.Errorf(
			"filter failed: %w",
			err,
		)
	}

	previous, err := FindPrevious(
		cfg.Output.Filter.Directory,
		cfg.Output.Filter.Suffix,
		filteredFile,
	)

	if err != nil {
		return err
	}

	//
	// 4. Generate diff
	//

	if previous != "" {

		if err := os.MkdirAll(
			cfg.Output.Diff.Directory,
			0755,
		); err != nil {
			return err
		}

		diffFile := filepath.Join(
			cfg.Output.Diff.Directory,
			fmt.Sprintf(
				"%s-%s",
				timestamp,
				cfg.Output.Diff.Suffix,
			),
		)

		fmt.Println(
			"Generating diff:",
			diffFile,
		)

		if err := diff.Generate(
			previous,
			filteredFile,
			diffFile,
		); err != nil {

			return fmt.Errorf(
				"diff generation failed: %w",
				err,
			)
		}
	}

	Cleanup(
		cfg.Output.Filter.Directory,
		cfg.Output.Filter.Suffix,
		cfg.Output.Filter.Retention,
	)

	Cleanup(
		cfg.Output.Diff.Directory,
		cfg.Output.Diff.Suffix,
		cfg.Output.Diff.Retention,
	)

	return nil
}
