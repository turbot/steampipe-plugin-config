---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/config.svg"
brand_color: "#222017"
display_name: "Config"
short_name: "config"
description: "Steampipe plugin to query data from various types of files, e.g. `.ini`, `.yml`, `.json`."
og_description: "Query data from various types of files with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/config-social-graphic.png"
---

# Config + Steampipe

Config plugin is used to parse various types of configuration files, e.g., `INI`, `JSON`, `YML`, in order to represent the content as SQL tables.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query data using SQL.

Query all data in your INI files:

```sql
select
  path,
  section,
  key,
  value
from
  ini_key_value;
```

```sh
> select section, key, value from ini_key_value where section = 'Settings';
+----------+---------------+-------------------------------------------+
| section  | key           | value                                     |
+----------+---------------+-------------------------------------------+
| Settings | DetailedLog   | 1                                         |
| Settings | RunStatus     | 1                                         |
| Settings | StatusRefresh | 10                                        |
| Settings | StatusPort    | 6090                                      |
| Settings | Archive       | 1                                         |
| Settings | ServerName    | Unknown                                   |
| Settings | LogFile       | /opt/ecs/mvuser/MV_IPTel/log/MV_IPTel.log |
| Settings | Version       | 0.9 Build 4 Created July 11 2004 14:00    |
+----------+---------------+-------------------------------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/config/tables)**

## Get started

### Install

Download and install the latest Config plugin:

```bash
steampipe plugin install config
```

### Credentials

No credentials are required.

### Configuration

Installing the latest config plugin will create a config file (`~/.steampipe/config/config.spc`) with a single connection named `config`:

```hcl
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
```

- `ini_paths` - A list of directory paths to search for INI files.
- `json_paths` - A list of directory paths to search for JSON files.
- `yml_paths` - A list of directory paths to search for YML files.

All above mentioned paths are resolved relative to the current working directory. Paths may [include wildcards](https://pkg.go.dev/path/filepath#Match) and also support `**` for recursive matching. Defaults to the current working directory.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-config
- Community: [Slack Channel](https://steampipe.io/community/join)
