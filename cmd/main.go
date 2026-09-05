package main

import (
	"flag"
	"fmt"
	"log"

	"ldap-extractor/config"
	"ldap-extractor/ldap"
	"ldap-extractor/version"
)

func main() {
	configFile := flag.String(
		"config",
		"config.yaml",
		"Path to configuration file",
	)

	outputFile := flag.String(
		"output",
		"",
		"Path to LDAP dump output file",
	)

	showVersion := flag.Bool(
		"version",
		false,
		"Show version information",
	)

	flag.Parse()

	if *showVersion {
		fmt.Println(version.String())
		return
	}

	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if *outputFile != "" {
		cfg.Output.File = *outputFile
	}

	if err := ldap.Dump(cfg); err != nil {
		log.Fatal(err)
	}
}
