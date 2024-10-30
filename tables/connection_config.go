package config

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type parseConfig struct {
	INIPaths  []string `hcl:"ini_paths,optional" steampipe:"watch"`
	JSONPaths []string `hcl:"json_paths,optional" steampipe:"watch"`
	TOMLPaths []string `hcl:"toml_paths,optional" steampipe:"watch"`
	YMLPaths  []string `hcl:"yml_paths,optional" steampipe:"watch"`
}

func ConfigInstance() interface{} {
	return &parseConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) parseConfig {
	if connection == nil || connection.Config == nil {
		return parseConfig{}
	}
	config, _ := connection.Config.(parseConfig)
	return config
}
