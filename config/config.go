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
	URL      string       `yaml:"url"`
	Username string       `yaml:"username"`
	Password string       `yaml:"password"`
	BaseDN   string       `yaml:"base_dn"`
	Search   SearchConfig `yaml:"search"`
	LDIF     LDIFConfig   `yaml:"ldif"`
}

type SearchConfig struct {
	Filter     string   `yaml:"filter"`
	PageSize   uint32   `yaml:"page_size"`
	Attributes []string `yaml:"attributes"`
}

type OutputConfig struct {
	File string `yaml:"file"`
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

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
