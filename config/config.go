package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LDAP   LDAPConfig   `yaml:"ldap"`
	Output OutputConfig `yaml:"output"`
}

type LDAPConfig struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	BaseDN   string `yaml:"base_dn"`

	Search SearchConfig `yaml:"search"`

	LDIF LDIFConfig `yaml:"ldif"`
}

type SearchConfig struct {
	Filter     string   `yaml:"filter"`
	PageSize   uint32   `yaml:"page_size"`
	Attributes []string `yaml:"attributes"`
}

type OutputConfig struct {

	// Raw LDAP dump
	File string `yaml:"file"`

	// Filtered JSON snapshots
	Filter FilterOutputConfig `yaml:"filter"`

	// Diff reports
	Diff DiffOutputConfig `yaml:"diff"`
}

type FilterOutputConfig struct {

	// Directory to store filtered JSON
	Directory string `yaml:"directory"`

	// Filename suffix
	// 20260909_120000-filtered.json
	Suffix string `yaml:"suffix"`

	// Number of filtered snapshots to retain
	Retention int `yaml:"retention"`
}

type DiffOutputConfig struct {

	// Directory to store diff files
	Directory string `yaml:"directory"`

	// Filename suffix
	// 20260909_120000-diff.json
	Suffix string `yaml:"suffix"`

	// Number of diff files to retain
	Retention int `yaml:"retention"`
}

type LDIFConfig struct {
	EncodeNonASCII bool `yaml:"encode_non_ascii"`
}

func Load(filename string) (*Config, error) {

	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := yaml.Unmarshal(
		data,
		&cfg,
	); err != nil {
		return nil, err
	}

	return &cfg, nil
}
