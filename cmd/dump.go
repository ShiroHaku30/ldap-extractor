package main

import (
	"github.com/spf13/cobra"

	"ldap-extractor/config"
	"ldap-extractor/ldap"
)

var (
	configFile string
	outputFile string
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump LDAP objects",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configFile)
		if err != nil {
			return err
		}

		if outputFile != "" {
			cfg.Output.File = outputFile
		}

		return ldap.Dump(cfg)
	},
}

func init() {
	dumpCmd.Flags().StringVarP(
		&configFile,
		"config",
		"c",
		"config.yaml",
		"Config file",
	)

	dumpCmd.Flags().StringVarP(
		&outputFile,
		"output",
		"o",
		"",
		"Output LDIF file",
	)

	rootCmd.AddCommand(dumpCmd)
}
