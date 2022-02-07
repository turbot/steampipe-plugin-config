connection "config" {
  plugin = "config"
  
  # Paths is a list of locations to search for files.
  # Wildcard based searches are supported, i.e. `/path/to/dir/*`, `*`, `**`, etc.
  # Exact file paths can have any name.
  paths = [ "./*" ]
}