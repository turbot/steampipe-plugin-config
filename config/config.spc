connection "config" {
  plugin = "config"
  
  # Paths is a list of locations to search for files.
  # Wildcard based searches are supported.
  # Exact file paths can have any name. Wildcard based matches must have an
  # extension, i.e. `.ini` (case insensitive).
  paths = [ "./*" ]
}