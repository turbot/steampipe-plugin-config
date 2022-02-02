---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/config.svg"
brand_color: "#808080"
display_name: "Config"
short_name: "config"
description: "Steampipe plugin to query data from various types of files, e.g. `.ini`."
og_description: "Query data from various types of files with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/config-social-graphic.png"
---

# Config + Steampipe

Config plugin is used to parse various types of files to represent the content as a SQL table.

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
  
  # Paths is a list of locations to search for files.
  # Wildcard based searches are supported.
  # Exact file paths can have any name. Wildcard based matches must have an
  # extension, i.e. `.ini` (case insensitive)(case insensitive).
  paths = [ "./*" ]
}
```

- `paths` - A list of directory paths to search for files. Paths may [include wildcards](https://pkg.go.dev/path/filepath#Match). File matches must have the required extension, i.e. `.ini` (case insensitive).

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-config
- Community: [Slack Channel](https://steampipe.io/community/join)
