package config

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"gopkg.in/yaml.v3"
)

func tableYMLKeyValue(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "yml_key_value",
		Description: "List all key value pairs from given YML file.",
		List: &plugin.ListConfig{
			Hydrate: listYMLKeyValue,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "path",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "Specifies the path of the YML file."},
			{Name: "key_path", Type: proto.ColumnType_LTREE, Transform: transform.FromField("Key").Transform(keysToSnakeCase), Description: "Specifies full path of a key in YML file."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "Specifies the value of the corresponding key."},
			{Name: "keys", Type: proto.ColumnType_JSON, Transform: transform.FromField("Key"), Description: "The array representation of path of a key."},
			{Name: "tag", Type: proto.ColumnType_STRING, Description: "Specifies the data type of the value."},
			{Name: "start_line", Type: proto.ColumnType_INT, Description: "Specifies the line number where the value is located."},
			{Name: "start_column", Type: proto.ColumnType_INT, Description: "Specifies the starting column of the value."},
			{Name: "pre_comments", Type: proto.ColumnType_JSON, Description: "Specifies the comments added above a key."},
			{Name: "head_comment", Type: proto.ColumnType_STRING, Description: "Specifies the comments in the lines preceding the node and not separated by an empty line."},
			{Name: "line_comment", Type: proto.ColumnType_STRING, Description: "Specifies the comments at the end of the line where the node is in."},
			{Name: "foot_comment", Type: proto.ColumnType_STRING, Description: "Specifies the comments following the node and before empty lines."},
		},
	}
}

func listYMLKeyValue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
		paths, err = fileList(ctx, d.Connection, ".yml")
		if err != nil {
			return nil, err
		}
	}

	for _, path := range paths {
		// Read file content
		reader, err := os.Open(path)
		if err != nil {
			// Could not open the file, so log and ignore
			plugin.Logger(ctx).Error("listYMLKeyValue", "file_error", err, "path", path)
			return nil, nil
		}

		var root yaml.Node
		decoder := yaml.NewDecoder(reader)
		err = decoder.Decode(&root)
		if err != nil {
			plugin.Logger(ctx).Error("listYMLKeyValue", "parse_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse file %s: %v", path, err)
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

type Rows []Row
type Row struct {
	Path        string
	Key         []string
	Value       interface{}
	Tag         *string
	PreComments []string
	HeadComment string
	LineComment string
	FootComment string
	StartLine   int
	StartColumn int
}

func treeToList(tree *yaml.Node, prefix []string, rows *Rows, preComments []string) {
	switch tree.Kind {
	case yaml.DocumentNode:
		for i, v := range tree.Content {
			localComments := []string{}
			if i == 0 {
				localComments = append(localComments, preComments...)
				if tree.HeadComment != "" {
					localComments = append(localComments, tree.HeadComment)
				}
				if tree.LineComment != "" {
					localComments = append(localComments, tree.LineComment)
				}
			}
			treeToList(v, prefix, rows, localComments)
		}
	case yaml.SequenceNode:
		if len(tree.Content) == 0 {
			row := Row{
				Key:         prefix,
				Value:       []string{},
				Tag:         tagToJSONType(tree.Tag),
				StartLine:   tree.Line,
				StartColumn: tree.Column,
				PreComments: preComments,
				HeadComment: tree.HeadComment,
				LineComment: tree.LineComment,
				FootComment: tree.FootComment,
			}
			*rows = append(*rows, row)
		}

		for i, v := range tree.Content {
			localComments := []string{}
			if i == 0 {
				localComments = append(localComments, preComments...)
				if tree.HeadComment != "" {
					localComments = append(localComments, tree.HeadComment)
				}
				if tree.LineComment != "" {
					localComments = append(localComments, tree.LineComment)
				}
			}
			newKey := make([]string, len(prefix))
			copy(newKey, prefix)
			newKey = append(newKey, strconv.Itoa(i))
			treeToList(v, newKey, rows, localComments)
		}
	case yaml.MappingNode:
		localComments := []string{}
		localComments = append(localComments, preComments...)
		if tree.HeadComment != "" {
			localComments = append(localComments, tree.HeadComment)
		}
		if tree.LineComment != "" {
			localComments = append(localComments, tree.LineComment)
		}
		if len(tree.Content) == 0 {
			row := Row{
				Key:         prefix,
				Value:       map[string]interface{}{},
				Tag:         tagToJSONType(tree.Tag),
				StartLine:   tree.Line,
				StartColumn: tree.Column,
				PreComments: preComments,
				HeadComment: tree.HeadComment,
				LineComment: tree.LineComment,
				FootComment: tree.FootComment,
			}
			*rows = append(*rows, row)
		}
		i := 0
		for i < len(tree.Content)-1 {
			key := tree.Content[i]
			val := tree.Content[i+1]
			i = i + 2
			if key.HeadComment != "" {
				localComments = append(localComments, key.HeadComment)
			}
			if key.LineComment != "" {
				localComments = append(localComments, key.LineComment)
			}
			newKey := make([]string, len(prefix))
			copy(newKey, prefix)
			newKey = append(newKey, key.Value)
			treeToList(val, newKey, rows, localComments)
			localComments = make([]string, 0)
		}
	case yaml.ScalarNode:
		row := Row{
			Key:         prefix,
			Value:       tree.Value,
			Tag:         tagToJSONType(tree.Tag),
			StartLine:   tree.Line,
			StartColumn: tree.Column,
			PreComments: preComments,
			HeadComment: tree.HeadComment,
			LineComment: tree.LineComment,
			FootComment: tree.FootComment,
		}
		if tree.Tag == "!!null" {
			row.Value = nil
		}
		*rows = append(*rows, row)
	}
}

func keysToSnakeCase(_ context.Context, d *transform.TransformData) (interface{}, error) {
	keys := d.Value.([]string)
	snakes := []string{}
	re := regexp.MustCompile(`[^A-Za-z0-9_]`)
	for _, k := range keys {
		snakes = append(snakes, re.ReplaceAllString(k, "_"))
	}
	return strings.Join(snakes, "."), nil
}

func tagToJSONType(tag string) *string {
	dataType := "string"
	switch tag {
	case "!!str":
		dataType = "string"
	case "!!int":
		dataType = "integer"
	case "!!float":
		dataType = "number"
	case "!!null":
		return nil
	case "!!bool":
		dataType = "boolean"
	}
	return &dataType
}
