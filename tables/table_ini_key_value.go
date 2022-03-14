package config

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tableINIKeyValue(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ini_key_value",
		Description: "Table representation of an INI file.",
		List: &plugin.ListConfig{
			Hydrate: listINIWithPath,
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
			{Name: "key", Type: proto.ColumnType_STRING, Description: "The name of the key."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "The value of corresponding key."},
			{Name: "comment", Type: proto.ColumnType_STRING, Description: "The short notes used to describe the key."},
		},
	}
}

type parseFormat struct {
	Path    string
	Section string
	Key     string
	Value   string
	Comment string
}

func listINIWithPath(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
		paths, err = listINIFiles(ctx, d.Connection)
		if err != nil {
			return nil, err
		}
	}

	for _, path := range paths {
		// Load file
		var opts ini.LoadOptions
		opts.AllowPythonMultilineValues = true
		cfg, err := ini.LoadSources(opts, path)
		if err != nil {
			plugin.Logger(ctx).Error("ini_key_value.listINIWithPath", "file_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse file %s: %v", path, err)
		}

		for _, i := range cfg.Sections() {
			// Extract keys of a section
			for _, key := range cfg.Section(i.Name()).Keys() {
				d.StreamListItem(ctx, formatResult(cfg, path, i.Name(), key.Name(), key.String(), key.Comment))
			}
		}
	}
	return nil, nil
}

func formatResult(cfg *ini.File, filePath string, secton string, key string, val string, comment string) parseFormat {
	return parseFormat{
		Path:    filePath,
		Section: secton,
		Key:     key,
		Value:   parseValue(cfg, val),
		Comment: comment,
	}
}

// parseValue will parse env variable and other variable references with its actual value
func parseValue(cfg *ini.File, str string) string {
	// Check for value of the environment variable references
	isEnvVar, _ := regexp.MatchString(".*\\${.*}.*", str)
	if isEnvVar {
		regexExp := regexp.MustCompile(`\$\{(.*?)\}`)
		matchedStr := regexExp.FindStringSubmatch(str)
		if len(matchedStr) > 1 {
			// Check for reference from other section, i.e. path = ${Common.system_dir}/Library/Frameworks/
			if strings.Contains(matchedStr[1], ".") {
				splitStr := strings.Split(matchedStr[1], ".")
				sec := strings.Join(splitStr[:len(splitStr)-1], ".")
				key := splitStr[len(splitStr)-1]
				value := cfg.Section(sec).Key(key).String()
				str = strings.Replace(str, matchedStr[0], value, -1)
			} else {
				// Replace the matched string with env variable value
				str = strings.Replace(str, matchedStr[0], os.Getenv(matchedStr[1]), -1)
			}
		}
	}
	return str
}
