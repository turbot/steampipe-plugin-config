package config

import (
	"context"
	"errors"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func listINIFiles(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	// Glob paths in config
	// Fail if no paths are specified
	parseConfig := GetConfig(d.Connection)
	if parseConfig.INIPaths == nil {
		return nil, errors.New("ini_paths must be configured to query INI files")
	}

	iniFiles, err := listPathsByFileType(ctx, d, parseConfig.INIPaths)
	if err != nil {
		return nil, err
	}
	return iniFiles, nil
}

func listYMLFiles(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	// Glob paths in config
	// Fail if no paths are specified
	parseConfig := GetConfig(d.Connection)
	if parseConfig.YMLPaths == nil {
		return nil, errors.New("yml_paths must be configured to query YML files")
	}

	ymlFiles, err := listPathsByFileType(ctx, d, parseConfig.YMLPaths)
	if err != nil {
		return nil, err
	}
	return ymlFiles, nil
}

func listJSONFiles(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	// Glob paths in config
	// Fail if no paths are specified
	parseConfig := GetConfig(d.Connection)
	if parseConfig.JSONPaths == nil {
		return nil, errors.New("json_paths must be configured to query JSON files")
	}

	jsonFiles, err := listPathsByFileType(ctx, d, parseConfig.JSONPaths)
	if err != nil {
		return nil, err
	}
	return jsonFiles, nil
}

func listPathsByFileType(ctx context.Context, d *plugin.QueryData, paths []string) ([]string, error) {
	var matches []string
	for _, i := range paths {
		// List the files in the given source directory
		files, err := d.GetSourceFiles(i)
		if err != nil {
			plugin.Logger(ctx).Error("utils.listPathsByFileType", "error getting source file info", err, "path", i)
			return nil, err
		}
		matches = append(matches, files...)
	}

	// Sanitize the matches to likely cloudformation files
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

		hit := false
		for _, j := range paths {
			if i == j {
				hit = true
				break
			}
		}
		if hit {
			fileList = append(fileList, i)
			continue
		}
		fileList = append(fileList, i)
	}
	return fileList, nil
}
