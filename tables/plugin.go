package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
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
			"ini_key_value": tableINIKeyValue(ctx),
			"ini_section":   tableINISection(ctx),
			"json_file":     tableJSONFile(ctx),
			"yml_file":      tableYMLFile(ctx),
			"yml_key_value": tableYMLKeyValue(ctx),
		},
	}
	return p
}

func fileList(ctx context.Context, p *plugin.Connection, fileType string) ([]string, error) {
	// Glob paths in config
	// Fail if no paths are specified
	parseConfig := GetConfig(p)
	if parseConfig.Paths == nil {
		return nil, errors.New("paths must be configured")
	}

	// Gather file path matches for the glob
	var matches []string
	paths := parseConfig.Paths
	for _, i := range paths {
		// Check to resolve ~ to home dir
		if strings.HasPrefix(i, "~") {
			// File system context
			home, err := os.UserHomeDir()
			if err != nil {
				plugin.Logger(ctx).Error("fileList", "os.UserHomeDir error. ~ will not be expanded in paths.", err)
			}

			// Resolve ~ to home dir
			if home != "" {
				if i == "~" {
					i = home
				} else if strings.HasPrefix(i, "~/") {
					i = filepath.Join(home, i[2:])
				}
			}
		}

		// Get full path
		fullPath, err := filepath.Abs(i)
		if err != nil {
			plugin.Logger(ctx).Error("fileList", "failed to fetch absolute path", err)
			return nil, err
		}

		if strings.Contains(fullPath, "*") {
			// Expand globs
			iMatches, err := doublestar.Glob(fullPath)
			if err != nil {
				return matches, fmt.Errorf("path is not a valid glob: %s", i)
			}
			matches = append(matches, iMatches...)
		} else {
			// Check file or diretory
			fileInfo, err := os.Stat(fullPath)
			if err != nil {
				plugin.Logger(ctx).Error("fileList", "error reading file path", err)
				return nil, err
			}

			// If directory, fetch all files
			if fileInfo.IsDir() {
				fullPath = filepath.Join(fullPath, "*")
			}

			// Expand globs
			iMatches, err := doublestar.Glob(fullPath)
			if err != nil {
				return matches, fmt.Errorf("path is not a valid glob: %s", fullPath)
			}
			matches = append(matches, iMatches...)
		}
	}

	// Sanitize the matches to likely files
	var filePaths []string
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
		switch fileType {
		case ".ini":
			if ext == ".ini" {
				filePaths = append(filePaths, i)
			}
		case ".json":
			if ext == ".json" {
				filePaths = append(filePaths, i)
			}
		case ".yml":
			if ext == ".yml" || ext == ".yaml" {
				filePaths = append(filePaths, i)
			}
		}
	}

	return filePaths, nil
}
