package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
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
			{Name: "key", Type: proto.ColumnType_LTREE, Description: "Specifies full path of a key in JSON file."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "Specifies the value of the corresponding key."},
			{Name: "keys", Type: proto.ColumnType_JSON, Description: "The array representation of path of a key."},
			{Name: "tag", Type: proto.ColumnType_STRING, Description: "Specifies the data type of the value."},
		},
	}
}

type parseJSONFormat struct {
	Path  string
	Key   string
	Value interface{}
	Keys  []string
	Tag   *string
}

func listJSONKeyValue(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Search for INI files to create as tables
	paths, err := fileList(ctx, d.Connection, ".json")
	if err != nil {
		return nil, err
	}

	if d.KeyColumnQuals["path"] != nil {
		paths = []string{d.KeyColumnQuals["path"].GetStringValue()}
	}

	for _, path := range paths {
		// Read file
		jsonFile, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("fail to read file: %v", err)
		}

		// defer the closing of jsonFile so that we can parse it later on
		defer jsonFile.Close()

		data := make(map[string]interface{})
		byteValue, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal([]byte(byteValue), &data)
		if err != nil {
			return nil, fmt.Errorf("fail to unmarshal data: %v", err)
		}

		result := map[string]interface{}{}
		for k, v := range data {
			if IsArray(v) {
				data := getNestedValue(k, v.([]interface{}), result)
				for k, v := range data {
					result[k] = v
				}
			} else if IsMap(v) {
				for mapKey, mapValue := range v.(map[string]interface{}) {
					newKey := fmt.Sprintf("%s.%s", k, mapKey)
					getNestedValue(newKey, mapValue, result)
				}
			} else {
				result[k] = v
			}
		}

		for k, v := range result {
			d.StreamListItem(ctx, parseJSONFormat{
				Path:  path,
				Key:   k,
				Value: v,
				Tag:   ValueToType(v),
				Keys:  strings.Split(k, "."),
			})
		}
	}
	return nil, nil
}

func getNestedValue(key string, value interface{}, result map[string]interface{}) map[string]interface{} {
	if IsArray(value) { // check array
		for count, i := range value.([]interface{}) {
			if IsMap(i) { // check array<map>
				for k, v := range i.(map[string]interface{}) {
					newKey := fmt.Sprintf("%s.%d.%s", key, count, k)
					getNestedValue(newKey, v, result)
				}
			} else {
				newKey := fmt.Sprintf("%s.%d", key, count)
				result[newKey] = i
			}
		}
	} else if IsMap(value) { // check map
		for k, v := range value.(map[string]interface{}) {
			newKey := fmt.Sprintf("%s.%s", key, k)
			getNestedValue(newKey, v, result)
		}
	} else {
		result[key] = value
	}

	return result
}
