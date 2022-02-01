package config

import "path/filepath"

func expandGlobs(path string) ([]string, error) {
	iMatches, err := filepath.Glob(path)
	if err != nil {
		// Fail if any path is an invalid glob
		return nil, err
	}
	return iMatches, nil
}
