package main

import (
	"github.com/spf13/cobra"

	"ldap-extractor/filter"
)

var (
	filterInputFile string
	fileOutputFile  string
)

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Extract selected LDAP attributes to JSON",
	RunE: func(cmd *cobra.Command, args []string) error {

		attributes := []string{
			"sAMAccountName",
			"department",
			"employeeID",
		}

		return filter.ExportJSON(
			filterInputFile,
			fileOutputFile,
			attributes,
		)
	},
}

func init() {

	filterCmd.Flags().StringVarP(
		&filterInputFile,
		"input",
		"i",
		"ldap_dump.ldif",
		"Input LDIF file",
	)

	filterCmd.Flags().StringVarP(
		&fileOutputFile,
		"output",
		"o",
		"users.json",
		"Output JSON file",
	)

	rootCmd.AddCommand(filterCmd)
}
