connection "config" {
  plugin = "config"
  
  # Paths is a list of locations to search for files.
  # Wildcards are supported per https://golang.org/pkg/path/filepath/#Match
  # Exact file paths can have any name. Wildcard based matches must have an
  # extension (case insensitive).
  paths = [ "./*" ]
}