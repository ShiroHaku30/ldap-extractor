package main

import (
	"github.com/spf13/cobra"

	"ldap-extractor/config"
	"ldap-extractor/extract"
)

var extractConfig string

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Dump LDAP, filter data, and generate diff",

	RunE: func(cmd *cobra.Command, args []string) error {

		cfg, err := config.Load(
			extractConfig,
		)

		if err != nil {
			return err
		}

		return extract.Run(cfg)
	},
}

func init() {

	extractCmd.Flags().
		StringVarP(
			&extractConfig,
			"config",
			"c",
			"config.yaml",
			"LDAP configuration file",
		)

	rootCmd.AddCommand(extractCmd)
}
