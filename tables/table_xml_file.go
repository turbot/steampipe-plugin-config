package config

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/clbanning/mxj/v2"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableXMLFile(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "xml_file",
		Description: "Represents the XML file content.",
		List: &plugin.ListConfig{
			Hydrate: listXMLFileWithPath,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "path",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "Specifies the path of the XML file."},
			{Name: "content", Type: proto.ColumnType_JSON, Description: "Specifies the file content."},
		},
	}
}

type parseXMLContent struct {
	Path    string
	Content interface{}
}

func listXMLFileWithPath(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var paths []string
	if d.EqualsQuals["path"] != nil {
		paths = []string{d.EqualsQuals["path"].GetStringValue()}
	} else {
		var err error
		paths, err = listXMLFiles(ctx, d)
		if err != nil {
			return nil, err
		}
	}

	for _, path := range paths {
		xmlFile, err := os.Open(path)
		if err != nil {
			plugin.Logger(ctx).Error("xml_file.listXMLFileWithPath", "file_error", err, "path", path)
			return nil, fmt.Errorf("fail to read file %s: %v", path, err)
		}

		byteValue, err := io.ReadAll(xmlFile)
		if cerr := xmlFile.Close(); cerr != nil {
			plugin.Logger(ctx).Error("xml_file.listXMLFileWithPath", "close_error", cerr, "path", path)
		}
		if err != nil {
			plugin.Logger(ctx).Error("xml_file.listXMLFileWithPath", "read_error", err, "path", path)
			return nil, fmt.Errorf("failed to read file content %s: %v", path, err)
		}

		mv, err := mxj.NewMapXml(byteValue)
		if err != nil {
			plugin.Logger(ctx).Error("xml_file.listXMLFileWithPath", "parse_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse XML content %s: %v", path, err)
		}

		d.StreamListItem(ctx, parseXMLContent{path, mv})
	}
	return nil, nil
}
