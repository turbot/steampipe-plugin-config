package config

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type parseConfig struct {
	INIPaths  []string `cty:"ini_paths"`
	JSONPaths []string `cty:"json_paths"`
	YMLPaths  []string `cty:"yml_paths"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"ini_paths": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
	"json_paths": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
	"yml_paths": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
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
