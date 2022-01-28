package config

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type parseConfig struct {
	Paths []string `cty:"paths"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"paths": {
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
