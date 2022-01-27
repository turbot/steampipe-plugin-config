package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-config",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromCamel().NullIfZero(),
		SchemaMode:       plugin.SchemaModeDynamic,
		TableMap: map[string]*plugin.Table{
			"config_ini_key_value": tableConfigINIKeyValue(ctx),
			"config_ini_section":   tableConfigINISection(ctx),
			"config_yaml":          tableConfigYAML(ctx),
		},
	}
	return p
}

func fileList(ctx context.Context, p *plugin.Connection, fileType string) ([]string, error) {

	var filePaths []string

	// Glob paths in config
	// Fail if no paths are specified
	parseConfig := GetConfig(p)
	if parseConfig.Paths == nil {
		return filePaths, errors.New("paths must be configured")
	}

	// File system context
	home, err := os.UserHomeDir()
	if err != nil {
		plugin.Logger(ctx).Error("fileList", "os.UserHomeDir error. ~ will not be expanded in paths.", err)
	}

	// Gather file path matches for the glob
	var matches []string
	paths := parseConfig.Paths
	for _, i := range paths {

		// Resolve ~ to home dir
		if home != "" {
			if i == "~" {
				i = home
			} else if strings.HasPrefix(i, "~/") {
				i = filepath.Join(home, i[2:])
			}
		}

		// Expand globs
		iMatches, err := filepath.Glob(i)
		if err != nil {
			// Fail if any path is an invalid glob
			return matches, fmt.Errorf("path is not a valid glob: %s", i)
		}

		matches = append(matches, iMatches...)
	}

	// Sanitize the matches to likely files
	for _, i := range matches {
		// If the file path is an exact match to a matrix path then it's always
		// treated as a match - it was requested exactly
		hit := false
		for _, j := range paths {
			if i == j {
				hit = true
				break
			}
		}
		if hit {
			filePaths = append(filePaths, i)
			continue
		}

		// This file was expanded from the glob, so check it's likely to be
		// of the right type based on the name / extension.
		ext := strings.ToLower(filepath.Ext(i))
		if ext == fileType {
			filePaths = append(filePaths, i)
		}
	}

	return filePaths, nil
}
