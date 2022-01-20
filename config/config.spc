connection "config" {
  plugin = "config"
  
  # Paths is a list of locations to search for files. Each file will be
  # converted to a table. Wildcards are supported per
  # https://golang.org/pkg/path/filepath/#Match
  # Exact file paths can have any name. Wildcard based matches must have an
  # extension of .ini (case insensitive).
  paths = [ "/path/to/dir/*", "/path/to/exact/custom.ini" ]
}