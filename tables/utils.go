package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func listINIFiles(ctx context.Context, p *plugin.Connection) ([]string, error) {
	// Glob paths in config
	// Fail if no paths are specified
	parseConfig := GetConfig(p)
	if parseConfig.INIPaths == nil {
		return nil, errors.New("ini_paths must be configured to query INI files")
	}

	iniFiles, err := listPathsByFileType(ctx, parseConfig.INIPaths)
	if err != nil {
		return nil, err
	}
	return iniFiles, nil
}

func listYMLFiles(ctx context.Context, p *plugin.Connection) ([]string, error) {
	// Glob paths in config
	// Fail if no paths are specified
	parseConfig := GetConfig(p)
	if parseConfig.YMLPaths == nil {
		return nil, errors.New("yml_paths must be configured to query YML files")
	}

	ymlFiles, err := listPathsByFileType(ctx, parseConfig.YMLPaths)
	if err != nil {
		return nil, err
	}
	return ymlFiles, nil
}

func listJSONFiles(ctx context.Context, p *plugin.Connection) ([]string, error) {
	// Glob paths in config
	// Fail if no paths are specified
	parseConfig := GetConfig(p)
	if parseConfig.JSONPaths == nil {
		return nil, errors.New("json_paths must be configured to query JSON files")
	}

	jsonFiles, err := listPathsByFileType(ctx, parseConfig.JSONPaths)
	if err != nil {
		return nil, err
	}
	return jsonFiles, nil
}

func listPathsByFileType(ctx context.Context, paths []string) ([]string, error) {
	var matches []string
	for _, i := range paths {
		// Check to resolve ~ to home dir
		if strings.HasPrefix(i, "~") {
			// File system context
			home, err := os.UserHomeDir()
			if err != nil {
				plugin.Logger(ctx).Error("utils.listPathsByFileType", "os.UserHomeDir error. ~ will not be expanded in paths", err, "path", i)
				return nil, fmt.Errorf("os.UserHomeDir error. ~ will not be expanded in paths")
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
			plugin.Logger(ctx).Error("utils.listPathsByFileType", "invlaid path", err, "path", i)
			return nil, fmt.Errorf("failed to fetch absolute path")
		}

		// Expand globs
		iMatches, err := doublestar.Glob(fullPath)
		if err != nil {
			plugin.Logger(ctx).Error("utils.listPathsByFileType", "path is not a valid glob", err, "path", i)
			return matches, fmt.Errorf("path is not a valid glob: %s", i)
		}
		matches = append(matches, iMatches...)
	}

	var fileList []string
	for _, i := range matches {
		// Check if file or directory
		fileInfo, err := os.Stat(i)
		if err != nil {
			plugin.Logger(ctx).Error("utils.listPathsByFileType", "error getting file info", err, "path", i)
			return nil, err
		}

		// Ignore directories
		if fileInfo.IsDir() {
			continue
		}
		fileList = append(fileList, i)
	}
	return fileList, nil
}
