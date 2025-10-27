package config

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableTOMLFile(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "toml_file",
		Description: "Represents the TOML file content.",
		List: &plugin.ListConfig{
			Hydrate: listTOMLFileWithPath,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "path",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "Specifies the path of the TOML file."},
			{Name: "content", Type: proto.ColumnType_JSON, Description: "Specifies the file content."},
		},
	}
}

type parseTOMLContent struct {
	Path    string
	Content interface{}
}

func listTOMLFileWithPath(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
		paths, err = listTOMLFiles(ctx, d)
		if err != nil {
			return nil, err
		}
	}

	for _, path := range paths {
		tomlFile, err := os.Open(path)
		if err != nil {
			plugin.Logger(ctx).Error("toml_file.listTOMLFileWithPath", "file_error", err, "path", path)
			return nil, fmt.Errorf("fail to read file %s: %v", path, err)
		}

		byteValue, err := io.ReadAll(tomlFile)
		if cerr := tomlFile.Close(); cerr != nil {
			plugin.Logger(ctx).Error("toml_file.listTOMLFileWithPath", "close_error", cerr, "path", path)
		}
		if err != nil {
			plugin.Logger(ctx).Error("toml_file.listTOMLFileWithPath", "read_error", err, "path", path)
			return nil, fmt.Errorf("failed to read file content %s: %v", path, err)
		}

		// Load TOML data
		var result interface{}
		err = toml.Unmarshal([]byte(byteValue), &result)
		if err != nil {
			plugin.Logger(ctx).Error("toml_file.listTOMLFileWithPath", "parse_error", err, "path", path)
			return nil, fmt.Errorf("failed to unmarshal file content %s: %v", path, err)
		}
		d.StreamListItem(ctx, parseTOMLContent{path, result})
	}
	return nil, nil
}
