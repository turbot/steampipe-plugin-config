connection "config" {
  plugin = "config"

  # Each paths argument is a list of locations to search for a particular file type
  # All paths are resolved relative to the current working directory (CWD)
  # Wildcard based searches are supported, including recursive searches.

  # For example, for the json_paths argument:
  #  - "*.json" matches all JSON files in the CWD
  #  - "**/*.json" matches all JSON files in a directory, and all the sub-directories in it
  #  - "../*.json" matches all JSON files in in the CWD's parent directory
  #  - "steampipe*.json" matches all JSON files starting with "steampipe" in the current CWD
  #  - "/path/to/dir/*.json" matches all JSON files in a specific directory
  #  - "/path/to/dir/main.json" matches a specific file

  # If any paths include "*", all files (including non-required files) in
  # the current CWD will be matched and will attempt to be loaded as that
  # particular file type

  # All paths arguments default to CWD
  ini_paths  = [ "*.ini" ]
  json_paths = [ "*.json" ]
  yml_paths  = [ "*.yml", "*.yaml" ]
}
