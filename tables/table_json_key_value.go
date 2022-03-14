package config

import (
	"context"
	"fmt"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"gopkg.in/yaml.v3"
)

func tableJSONKeyValue(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "json_key_value",
		Description: "List all key value pairs from given JSON file.",
		List: &plugin.ListConfig{
			Hydrate: listJSONKeyValue,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "path",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "Specifies the path of the JSON file."},
			{Name: "key_path", Type: proto.ColumnType_LTREE, Transform: transform.FromField("Key").Transform(keysToSnakeCase), Description: "Specifies full path of a key in JSON file."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "Specifies the value of the corresponding key."},
			{Name: "keys", Type: proto.ColumnType_JSON, Transform: transform.FromField("Key"), Description: "The array representation of path of a key."},
			{Name: "tag", Type: proto.ColumnType_STRING, Description: "Specifies the data type of the value."},
			{Name: "start_line", Type: proto.ColumnType_INT, Description: "Specifies the line number where the value is located."},
			{Name: "start_column", Type: proto.ColumnType_INT, Description: "Specifies the starting column of the value."},
		},
	}
}

func listJSONKeyValue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// #1 - Path via qual
	// If the path was requested through qualifier then match it exactly. Globs
	// are not supported in this context since the output value for the column
	// will never match the requested value.
	//
	// #2 - Path via glob paths in config
	var paths []string
	if d.KeyColumnQuals["path"] != nil {
		paths = []string{d.KeyColumnQuals["path"].GetStringValue()}
	} else {
		var err error
		paths, err = listJSONFiles(ctx, d.Connection)
		if err != nil {
			return nil, err
		}
	}

	for _, path := range paths {
		// Read file content
		reader, err := os.Open(path)
		if err != nil {
			// Could not open the file, so log and ignore
			plugin.Logger(ctx).Error("json_key_value.listJSONKeyValue", "file_error", err, "path", path)
			return nil, nil
		}

		var root yaml.Node
		decoder := yaml.NewDecoder(reader)
		err = decoder.Decode(&root)
		if err != nil {
			plugin.Logger(ctx).Error("json_key_value.listJSONKeyValue", "parse_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse file: %v", err)
		}

		var rows Rows
		treeToList(&root, []string{}, &rows, nil)
		for _, r := range rows {
			r.Path = path
			d.StreamListItem(ctx, r)
		}
	}
	return nil, nil
}
