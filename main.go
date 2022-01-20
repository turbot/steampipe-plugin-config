package main

import (
	config "github.com/turbot/steampipe-plugin-config/config_"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: config.Plugin})
}
