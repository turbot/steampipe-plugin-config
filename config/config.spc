connection "config" {
  plugin = "config"

  # Each paths argument is a list of locations to search for a particular file type
  # All paths are resolved relative to the current working directory (CWD)
  # Wildcard based searches are supported, including recursive searches.

  # All paths arguments default to CWD
  ini_paths  = [ "*.ini" ]
  json_paths = [ "*.json" ]
  yml_paths  = [ "*.yml", "*.yaml" ]
}
