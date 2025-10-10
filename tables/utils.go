package config

import (
	"context"
	"errors"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func listFilesByType(
	ctx context.Context,
	d *plugin.QueryData,
	paths []string,
	errMsg string,
) ([]string, error) {
	if paths == nil {
		return nil, errors.New(errMsg)
	}
	return listPathsByFileType(ctx, d, paths)
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

func listINIFiles(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	cfg := GetConfig(d.Connection)
	return listFilesByType(ctx, d, cfg.INIPaths, "ini_paths must be configured to query INI files")
}

func listXMLFiles(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	cfg := GetConfig(d.Connection)
	return listFilesByType(ctx, d, cfg.XMLPaths, "xml_paths must be configured to query XML files")
}

func listYMLFiles(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	cfg := GetConfig(d.Connection)
	return listFilesByType(ctx, d, cfg.YMLPaths, "yml_paths must be configured to query YML files")
}

func listJSONFiles(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	cfg := GetConfig(d.Connection)
	return listFilesByType(ctx, d, cfg.JSONPaths, "json_paths must be configured to query JSON files")
}

func listTOMLFiles(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	cfg := GetConfig(d.Connection)
	return listFilesByType(ctx, d, cfg.TOMLPaths, "toml_paths must be configured to query TOML files")
}
