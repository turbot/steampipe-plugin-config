package config

import (
	"context"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableConfigYAML(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "config_yaml",
		Description: "Table representation of an YAML file.",
		List: &plugin.ListConfig{
			Hydrate: listYAMLWithPath,
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "path",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "Specifies the path of the ini file."},
			{Name: "key", Type: proto.ColumnType_STRING, Description: "The name of the key."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "The value of corresponding key."},
		},
	}
}

type parseYamlFormat struct {
	Path  string
	Key   string
	Value interface{}
}

func listYAMLWithPath(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Search for INI files to create as tables
	paths, err := fileList(ctx, d.Connection, ".yaml")
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
		data := make(map[interface{}]interface{})
		err = yaml.Unmarshal(content, &data)
		if err != nil {
			return nil, fmt.Errorf("fail to unmarshal file content: %v", err)
		}

		result := map[string]interface{}{}
		for k, v := range data {
			r := map[string]interface{}{}
			content, isArray := v.([]interface{})
			if isArray {
				data := getNestedValue(k.(string), content, r)
				for k, v := range data {
					result[k] = v
				}
			} else {
				result[k.(string)] = v.(string)
			}
		}

		for k, v := range result {
			d.StreamListItem(ctx, parseYamlFormat{path, k, v})
		}
	}
	return nil, nil
}

// 1. First check is the content is array or not
// 2. If array
// 3.Check if it's a array[map]?
//   If array[map], continue step 1-4
//   else return the key as <parent>.<key>
// i.e. if array[map], then a[0].key1, a[0].key2
// if only map, the nested will be represented as a.b1.c1, a.b2.c2
// 4. else return the array of elements as it is
func getNestedValue(key string, value interface{}, result map[string]interface{}) map[string]interface{} {
	content, isArray := value.([]interface{})
	mapContent, isMap := value.(map[string]interface{})
	if isArray {
		// Check array length
		if len(content) > 0 {
			_, isMap := content[0].(map[string]interface{})
			if isMap {
				for count, i := range content {
					newKey := fmt.Sprintf("%s[%d]", key, count)
					getNestedValue(newKey, i, result)
				}
			} else {
				result[key] = content
			}
		}
	} else if isMap {
		for k, v := range mapContent {
			newKey := fmt.Sprintf("%s.%s", key, k)
			getNestedValue(newKey, v, result)
		}
	} else {
		result[key] = value
	}

	return result
}
