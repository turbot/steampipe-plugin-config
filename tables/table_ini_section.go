package config

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableINISection(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ini_section",
		Description: "Retrieves a list of sections and subsections defined in a INI file.",
		List: &plugin.ListConfig{
			Hydrate: listINISections,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "path",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "Specifies the path of the INI file."},
			{Name: "section", Type: proto.ColumnType_STRING, Description: "Specifies the name of the section."},
			{Name: "comment", Type: proto.ColumnType_STRING, Description: "The short notes used to describe the key."},
		},
	}
}

type parseSectionFormat struct {
	Path    string
	Section string
	Comment string
}

func listINISections(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// #1 - Path via qual
	// If the path was requested through qualifier then match it exactly. Globs
	// are not supported in this context since the output value for the column
	// will never match the requested value.
	//
	// #2 - Path via glob paths in config
	var paths []string
	if d.KeyColumnQuals["path"] != nil {
		ext := strings.ToLower(filepath.Ext(d.KeyColumnQuals["path"].GetStringValue()))
		if ext == ".ini" {
			paths = []string{d.KeyColumnQuals["path"].GetStringValue()}
		}
	} else {
		var err error
		paths, err = fileList(ctx, d.Connection, ".ini")
		if err != nil {
			return nil, err
		}
	}

	for _, path := range paths {
		// Load file
		var opts ini.LoadOptions
		cfg, err := ini.LoadSources(opts, path)
		if err != nil {
			plugin.Logger(ctx).Error("ini_section.listINISections", "file_error", err, "path", path)
			return nil, fmt.Errorf("failed to read file: %v", err)
		}

		for _, i := range cfg.Sections() {
			d.StreamListItem(ctx, parseSectionFormat{path, i.Name(), i.Comment})
		}
	}
	return nil, nil
}
