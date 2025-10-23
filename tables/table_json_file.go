package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableJSONFile(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "json_file",
		Description: "Represents the JSON file content.",
		List: &plugin.ListConfig{
			Hydrate: listJSONFileWithPath,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "path",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "Specifies the path of the JSON file."},
			{Name: "content", Type: proto.ColumnType_JSON, Description: "Specifies the file content."},
		},
	}
}

type parseJSONContent struct {
	Path    string
	Content interface{}
}

func listJSONFileWithPath(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// #1 - Path via qual
	// If the path was requested through qualifier then match it exactly. Globs
	// are not supported in this context since the output value for the column
	// will never match the requested value.
	//
	// #2 - Path via glob paths in config
	var paths []string
	if d.EqualsQuals["path"] != nil {
		paths = []string{d.EqualsQuals["path"].GetStringValue()}
	} else {
		var err error
		paths, err = listJSONFiles(ctx, d)
		if err != nil {
			return nil, err
		}
	}

	for _, path := range paths {
		jsonFile, err := os.Open(path)
		if err != nil {
			plugin.Logger(ctx).Error("json_file.listJSONFileWithPath", "file_error", err, "path", path)
			return nil, fmt.Errorf("fail to read file %s: %v", path, err)
		}

		byteValue, err := io.ReadAll(jsonFile)
		jsonFile.Close()
		if err != nil {
			plugin.Logger(ctx).Error("json_file.listJSONFileWithPath", "read_error", err, "path", path)
			return nil, fmt.Errorf("failed to read file content %s: %v", path, err)
		}

		// Load either JSON objects or JSON arrays
		var result interface{}
		err = json.Unmarshal([]byte(byteValue), &result)
		if err != nil {
			plugin.Logger(ctx).Error("json_file.listJSONFileWithPath", "parse_error", err, "path", path)
			return nil, fmt.Errorf("failed to unmarshal file content %s: %v", path, err)
		}
		d.StreamListItem(ctx, parseJSONContent{path, result})
	}
	return nil, nil
}
