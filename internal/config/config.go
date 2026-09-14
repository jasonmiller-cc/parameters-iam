// Package config provides configuration loading for the parameters-iam service.
package config

import (
	coreconfig "github.com/jasonmiller-cc/parameters-core/pkg/config"
)

// IAMConfig holds IAM-specific backend configuration.
type IAMConfig struct {
	LDAPServiceURL     string   `yaml:"ldap_service_url"     env:"IAM_LDAP_SERVICE_URL"`
	KerberosServiceURL string   `yaml:"kerberos_service_url" env:"IAM_KERBEROS_SERVICE_URL"`
	CAServiceURL       string   `yaml:"ca_service_url"       env:"IAM_CA_SERVICE_URL"`
	DefaultRoles       []string `yaml:"default_roles"`
	SuperAdminGroup    string   `yaml:"super_admin_group"    env:"IAM_SUPER_ADMIN_GROUP"`
}

// Config is the top-level configuration for parameters-iam.
type Config struct {
	coreconfig.BaseConfig `yaml:",inline"`
	IAM                   IAMConfig `yaml:"iam"`
}

// Load reads a YAML config file at path (empty string falls back to standard
// locations) and overlays environment variables with the PARAMS_IAM prefix.
func Load(path string) (*Config, error) {
	cfg := &Config{}
	if err := coreconfig.Load(path, "PARAMS_IAM", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
