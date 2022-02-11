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

The following file types are currently supported:
- INI
- JSON
- YML

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
+--------------------------------+----------+---------------+-------------------------------------------+
| path                           | section  | key           | value                                     |
+--------------------------------+----------+---------------+-------------------------------------------+
| /Users/myuser/ini/defaults.ini | Settings | DetailedLog   | 1                                         |
| /Users/myuser/ini/defaults.ini | Status   | RunStatus     | 1                                         |
| /Users/myuser/ini/defaults.ini | Status   | StatusRefresh | 10                                        |
| /Users/myuser/ini/defaults.ini | Status   | StatusPort    | 6090                                      |
| /Users/myuser/ini/logs.ini     | Server   | Archive       | 1                                         |
| /Users/myuser/ini/logs.ini     | Server   | ServerName    | Unknown                                   |
| /Users/myuser/ini/logs.ini     | Settings | LogFile       | /opt/ecs/mvuser/MV_IPTel/log/MV_IPTel.log |
| /Users/myuser/ini/logs.ini     | Settings | Version       | 0.9 Build 4 Created July 11 2004 14:00    |
+--------------------------------+----------+---------------+-------------------------------------------+
```

Query all data in your JSON files:

```sql
select
  path,
  jsonb_pretty(content) as file_content
from
  json_file;
```

```sh
+---------------------------------+------------------------------------------------------------------------------------------------------------------------------+
| path                            | file_content                                                                                                                 |
+---------------------------------+------------------------------------------------------------------------------------------------------------------------------+
| /Users/myuser/json/invoice.json | {                                                                                                                            |
|                                 |     "city": "East Centerville",                                                                                              |
|                                 |     "date": "2012-08-06T00:00:00Z",                                                                                          |
|                                 |     "items": [                                                                                                               |
|                                 |         {                                                                                                                    |
|                                 |             "price": 1.47,                                                                                                   |
|                                 |             "part_no": "A4786",                                                                                              |
|                                 |             "quantity": 4,                                                                                                   |
|                                 |             "description": "Water Bucket (Filled)"                                                                           |
|                                 |         },                                                                                                                   |
|                                 |         {                                                                                                                    |
|                                 |             "size": 8,                                                                                                       |
|                                 |             "price": 133.7,                                                                                                  |
|                                 |             "part_no": "E1628",                                                                                              |
|                                 |             "quantity": 1,                                                                                                   |
|                                 |             "description": "High Heeled \"Ruby\" Slippers"                                                                   |
|                                 |         }                                                                                                                    |
|                                 |     ],                                                                                                                       |
|                                 |     "state": "KS",                                                                                                           |
|                                 |     "street": "123 Tornado Alley\nSuite 16\n",                                                                               |
|                                 |     "bill-to": null,                                                                                                         |
|                                 |     "receipt": "Oz-Ware Purchase Invoice",                                                                                   |
|                                 |     "ship-to": null,                                                                                                         |
|                                 |     "customer": {                                                                                                            |
|                                 |         "first_name": "Dorothy",                                                                                             |
|                                 |         "family_name": "Gale"                                                                                                |
|                                 |     },                                                                                                                       |
|                                 |     "specialDelivery": "Follow the Yellow Brick Road to the Emerald City. Pay no attention to the man behind the curtain.\n" |
|                                 | }                                                                                                                            |
| /Users/myuser/json/test.json    | {                                                                                                                            |
|                                 |     "foo": "bar",                                                                                                            |
|                                 |     "includes": [                                                                                                            |
|                                 |         "common.json"                                                                                                        |
|                                 |     ]                                                                                                                        |
|                                 | }                                                                                                                            |
+---------------------------------+------------------------------------------------------------------------------------------------------------------------------+
```

Query all data in your YML files:

```sql
select
  path,
  jsonb_pretty(content) as file_content
from
  yml_file;
```

```sh
+--------------------------------+------------------------------------------------------------------------------------------------------------------------------+
| path                           | file_content                                                                                                                 |
+--------------------------------+------------------------------------------------------------------------------------------------------------------------------+
| /Users/myuser/yml/invoice.yml  | {                                                                                                                            |
|                                |     "city": "East Centerville",                                                                                              |
|                                |     "date": "2012-08-06T00:00:00Z",                                                                                          |
|                                |     "items": [                                                                                                               |
|                                |         {                                                                                                                    |
|                                |             "price": 1.47,                                                                                                   |
|                                |             "part_no": "A4786",                                                                                              |
|                                |             "quantity": 4,                                                                                                   |
|                                |             "description": "Water Bucket (Filled)"                                                                           |
|                                |         },                                                                                                                   |
|                                |         {                                                                                                                    |
|                                |             "size": 8,                                                                                                       |
|                                |             "price": 133.7,                                                                                                  |
|                                |             "part_no": "E1628",                                                                                              |
|                                |             "quantity": 1,                                                                                                   |
|                                |             "description": "High Heeled \"Ruby\" Slippers"                                                                   |
|                                |         }                                                                                                                    |
|                                |     ],                                                                                                                       |
|                                |     "state": "KS",                                                                                                           |
|                                |     "street": "123 Tornado Alley\nSuite 16\n",                                                                               |
|                                |     "bill-to": null,                                                                                                         |
|                                |     "receipt": "Oz-Ware Purchase Invoice",                                                                                   |
|                                |     "ship-to": null,                                                                                                         |
|                                |     "customer": {                                                                                                            |
|                                |         "first_name": "Dorothy",                                                                                             |
|                                |         "family_name": "Gale"                                                                                                |
|                                |     },                                                                                                                       |
|                                |     "specialDelivery": "Follow the Yellow Brick Road to the Emerald City. Pay no attention to the man behind the curtain.\n" |
|                                | }                                                                                                                            |
| /Users/myuser/yml/test.yaml    | {                                                                                                                            |
|                                |     "foo": "bar",                                                                                                            |
|                                |     "includes": [                                                                                                            |
|                                |         "common.yaml"                                                                                                        |
|                                |     ]                                                                                                                        |
|                                | }                                                                                                                            |
+--------------------------------+------------------------------------------------------------------------------------------------------------------------------+
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
```

- `ini_paths` - A list of directory paths to search for INI files.
- `json_paths` - A list of directory paths to search for JSON files.
- `yml_paths` - A list of directory paths to search for YML files.

All `paths` arguments are resolved relative to the current working directory. Paths may [include wildcards](https://pkg.go.dev/path/filepath#Match) and also support `**` for recursive matching. Each `paths` argument defaults to the current working directory.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-config
- Community: [Slack Channel](https://steampipe.io/community/join)
