package main

import (
	"github.com/spf13/cobra"

	"ldap-extractor/config"
	"ldap-extractor/ldap"
)

var testConfig string

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test LDAP connectivity",
	RunE: func(cmd *cobra.Command, args []string) error {

		cfg, err := config.Load(testConfig)
		if err != nil {
			return err
		}

		return ldap.TestConnection(cfg)
	},
}

func init() {

	testCmd.Flags().StringVarP(
		&testConfig,
		"config",
		"c",
		"config.yaml",
		"LDAP configuration file",
	)

	rootCmd.AddCommand(testCmd)
}
