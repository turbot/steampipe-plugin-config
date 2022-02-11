connection "config" {
  plugin = "config"

  # Paths is a list of locations to search for various types of files.
  # Configure paths based on your file type. For example:
  # Use "ini_paths" to search for INI files. Similarly, use "json_paths" and "yml_paths" for JSON and YML files respectively.
  # Wildcard based searches are supported, including recursive searches.

  # For example:
  #  - "*.json" matches all JSON files in the CWD
  #  - "**/*.json" matches all JSON files in a directory, and all the sub-directories in it
  #  - "../*.json" matches all JSON files in in the CWD's parent directory
  #  - "steampipe*.json" matches all JSON files starting with "steampipe" in the current CWD
  #  - "/path/to/dir/*.json" matches all JSON files in a specific directory
  #  - "/path/to/dir/main.json" matches a specific file

  # If paths includes "*", all files (including non-required files) in
  # the current CWD will be matched, which may cause errors if incompatible filetypes exist

  # Defaults to CWD
  ini_paths = [ "*.ini" ]
  json_paths = [ "*.json" ]
  yml_paths = [ "*.yml", "*.yaml" ]
}