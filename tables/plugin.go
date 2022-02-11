package config

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-config",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromCamel().NullIfZero(),
		SchemaMode:       plugin.SchemaModeDynamic,
		TableMap: map[string]*plugin.Table{
			"ini_key_value": tableINIKeyValue(ctx),
			"ini_section":   tableINISection(ctx),
			"json_file":     tableJSONFile(ctx),
			"yml_file":      tableYMLFile(ctx),
		},
	}
	return p
}
