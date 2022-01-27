package config

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"gopkg.in/yaml.v3"
)

func tableConfigYMLFile(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "config_yml_file",
		Description: "Represents the YML file content into JSON format.",
		List: &plugin.ListConfig{
			Hydrate: listYMLFileWithPath,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "path",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "Specifies the path of the ini file."},
			{Name: "content", Type: proto.ColumnType_JSON, Description: "Specifies the file content in JSON format."},
		},
	}
}

type parseYMLContent struct {
	Path    string
	Content interface{}
}

func listYMLFileWithPath(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Search for INI files to create as tables
	paths, err := fileList(ctx, d.Connection, ".yml")
	if err != nil {
		return nil, err
	}

	if d.KeyColumnQuals["path"] != nil {
		paths = []string{d.KeyColumnQuals["path"].GetStringValue()}
	}

	for _, path := range paths {
		// Read file
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("fail to read file: %v", err)
		}

		// Decoding the file content
		var data interface{}
		err = yaml.Unmarshal(content, &data)
		if err != nil {
			return nil, fmt.Errorf("fail to unmarshal file content: %v", err)
		}
		d.StreamListItem(ctx, parseYMLContent{path, data})
	}
	return nil, nil
}
