package loader

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LDAP struct {
		URL      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		BaseDN   string `yaml:"base_dn"`

		Search struct {
			Filter     string   `yaml:"filter"`
			PageSize   uint32   `yaml:"page_size"`
			Attributes []string `yaml:"attributes"`
		} `yaml:"search"`
	} `yaml:"ldap"`

	Output struct {
		File string `yaml:"file"`
	} `yaml:"output"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
