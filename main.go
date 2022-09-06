package main

import (
	config "github.com/turbot/steampipe-plugin-config/tables"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: config.Plugin})
}
