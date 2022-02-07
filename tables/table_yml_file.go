package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"gopkg.in/yaml.v3"
)

func tableYMLFile(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yml_file",
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
			{Name: "path", Type: proto.ColumnType_STRING, Description: "Specifies the path of the YML file."},
			{Name: "content", Type: proto.ColumnType_JSON, Description: "Specifies the file content in JSON format."},
		},
	}
}

type parseYMLContent struct {
	Path    string
	Content interface{}
}

func listYMLFileWithPath(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// #1 - Path via qual
	// If the path was requested through qualifier then match it exactly. Globs
	// are not supported in this context since the output value for the column
	// will never match the requested value.
	//
	// #2 - Path via glob paths in config
	var paths []string
	if d.KeyColumnQuals["path"] != nil {
		ext := strings.ToLower(filepath.Ext(d.KeyColumnQuals["path"].GetStringValue()))
		if ext == ".yml" || ext == ".yaml" {
			paths = []string{d.KeyColumnQuals["path"].GetStringValue()}
		}
	} else {
		var err error
		paths, err = fileList(ctx, d.Connection, ".yml")
		if err != nil {
			return nil, err
		}
	}

	for _, path := range paths {
		// Read file
		content, err := ioutil.ReadFile(path)
		if err != nil {
			plugin.Logger(ctx).Error("yml_file.listYMLFileWithPath", "file_error", err, "path", path)
			return nil, fmt.Errorf("failed to read file: %v", err)
		}

		// Decoding the file content
		var data interface{}
		err = yaml.Unmarshal(content, &data)
		if err != nil {
			plugin.Logger(ctx).Error("yml_file.listYMLFileWithPath", "parse_error", err, "path", path)
			return nil, fmt.Errorf("failed to unmarshal file content: %v", err)
		}
		d.StreamListItem(ctx, parseYMLContent{path, data})
	}
	return nil, nil
}
