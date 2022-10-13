---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/config.svg"
brand_color: "#222017"
display_name: "Config"
short_name: "config"
description: "Steampipe plugin to query data from various types of files like INI, JSON, YML and more."
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
+----------------------------+----------+---------------+-------------------------------------------+
| path                       | section  | key           | value                                     |
+----------------------------+----------+---------------+-------------------------------------------+
| /Users/myuser/defaults.ini | Settings | DetailedLog   | 1                                         |
| /Users/myuser/defaults.ini | Status   | RunStatus     | 1                                         |
| /Users/myuser/defaults.ini | Status   | StatusRefresh | 10                                        |
| /Users/myuser/defaults.ini | Status   | StatusPort    | 6090                                      |
| /Users/myuser/logs.ini     | Server   | Archive       | 1                                         |
| /Users/myuser/logs.ini     | Server   | ServerName    | Unknown                                   |
| /Users/myuser/logs.ini     | Settings | LogFile       | /opt/ecs/mvuser/MV_IPTel/log/MV_IPTel.log |
| /Users/myuser/logs.ini     | Settings | Version       | 0.9 Build 4 Created July 11 2004 14:00    |
+----------------------------+----------+---------------+-------------------------------------------+
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
+----------------------------+------------------------------------------------------------+
| path                       | file_content                                               |
+----------------------------+------------------------------------------------------------+
| /Users/myuser/invoice.json | {                                                          |
|                            |     "city": "East Centerville",                            |
|                            |     "date": "2012-08-06T00:00:00Z",                        |
|                            |     "items": [                                             |
|                            |         {                                                  |
|                            |             "price": 1.47,                                 |
|                            |             "part_no": "A4786",                            |
|                            |             "quantity": 4,                                 |
|                            |             "description": "Water Bucket (Filled)"         |
|                            |         },                                                 |
|                            |         {                                                  |
|                            |             "size": 8,                                     |
|                            |             "price": 133.7,                                |
|                            |             "part_no": "E1628",                            |
|                            |             "quantity": 1,                                 |
|                            |             "description": "High Heeled \"Ruby\" Slippers" |
|                            |         }                                                  |
|                            |     ],                                                     |
|                            |     "state": "KS",                                         |
|                            |     "street": "123 Tornado Alley\nSuite 16\n",             |
|                            |     "bill-to": null,                                       |
|                            |     "receipt": "Oz-Ware Purchase Invoice",                 |
|                            |     "ship-to": null,                                       |
|                            |     "customer": {                                          |
|                            |         "first_name": "Dorothy",                           |
|                            |         "family_name": "Gale"                              |
|                            |     },                                                     |
|                            | }                                                          |
| /Users/myuser/test.json    | {                                                          |
|                            |     "foo": "bar",                                          |
|                            |     "includes": [                                          |
|                            |         "common.json"                                      |
|                            |     ]                                                      |
|                            | }                                                          |
+----------------------------+------------------------------------------------------------+
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
+---------------------------+------------------------------------------------------------+
| path                      | file_content                                               |
+---------------------------+------------------------------------------------------------+
| /Users/myuser/invoice.yml | {                                                          |
|                           |     "city": "East Centerville",                            |
|                           |     "date": "2012-08-06T00:00:00Z",                        |
|                           |     "items": [                                             |
|                           |         {                                                  |
|                           |             "price": 1.47,                                 |
|                           |             "part_no": "A4786",                            |
|                           |             "quantity": 4,                                 |
|                           |             "description": "Water Bucket (Filled)"         |
|                           |         },                                                 |
|                           |         {                                                  |
|                           |             "size": 8,                                     |
|                           |             "price": 133.7,                                |
|                           |             "part_no": "E1628",                            |
|                           |             "quantity": 1,                                 |
|                           |             "description": "High Heeled \"Ruby\" Slippers" |
|                           |         }                                                  |
|                           |     ],                                                     |
|                           |     "state": "KS",                                         |
|                           |     "street": "123 Tornado Alley\nSuite 16\n",             |
|                           |     "bill-to": null,                                       |
|                           |     "receipt": "Oz-Ware Purchase Invoice",                 |
|                           |     "ship-to": null,                                       |
|                           |     "customer": {                                          |
|                           |         "first_name": "Dorothy",                           |
|                           |         "family_name": "Gale"                              |
|                           |     },                                                     |
|                           | }                                                          |
| /Users/myuser/test.yaml   | {                                                          |
|                           |     "foo": "bar",                                          |
|                           |     "includes": [                                          |
|                           |         "common.yaml"                                      |
|                           |     ]                                                      |
|                           | }                                                          |
+---------------------------+------------------------------------------------------------+
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

  # All paths arguments default to CWD
  ini_paths  = [ "*.ini" ]
  json_paths = [ "*.json" ]
  yml_paths  = [ "*.yml", "*.yaml" ]
}
```

- `ini_paths` - A list of directory paths to search for INI files.
- `json_paths` - A list of directory paths to search for JSON files.
- `yml_paths` - A list of directory paths to search for YML files.

### Setting up paths

The arguments `ini_paths`, `json_paths` and `yml_paths` in the config are a list of directory paths, a GitHub repository URL, or a S3 URL to search for INI, JSON and YML files respectively. Paths may [include wildcards](https://pkg.go.dev/path/filepath#Match) and also support `**` for recursive matching. Defaults to the current working directory.

#### Configuring local file paths

You can define a list of local directory paths to search for INI, JSON and YML files. Paths are resolved relative to the current working directory. For example, for the `json_paths` argument:

- `*.json` matches all JSON files in the CWD.
- `**/*.json` matches all JSON files in the CWD and all sub-directories.
- `../*.json` matches all JSON files in the CWD's parent directory.
- `steampipe*.json` matches all JSON files starting with `steampipe` in the CWD.
- `/path/to/dir/*.json` matches all JSON files in a specific directory. For example:
  - `~/*.json` matches all JSON files in the home directory.
  - `~/**/*.json` matches all JSON files recursively in the home directory.
- `/path/to/dir/main.json` matches a specific file.

**NOTE:** If paths includes `*`, all files (including non-required files) in the CWD will be matched, which may cause errors if incompatible file types exist.

#### Configuring GitHub URLs

You can define a list of URL as input to search for INI, JSON and YML files from a variety of protocols. For example:

- `github.com/LearnWebCode/json-example//*.json` matches all top-level JSON files in the specified github repository.
- `github.com/LearnWebCode/json-example//**/*.json` matches all JSON files in the specified github repository and all sub-directories.
- `github.com/turbot/polygoat?ref=fix_7677//**/*.json` matches all JSON files in the specific tag of github repository.

If you want to download only a specific subdirectory from a downloaded directory, you can specify a subdirectory after a double-slash (`//`).

- `github.com/brandiqa/json-examples//api/config//*.json` matches all JSON files in the specific folder of a github repository.

#### Configuring S3 URLs

You can also pass a S3 bucket URL to search for INI, JSON and YML files stored in the specified S3 bucket. For example:

- `s3::https://s3.amazonaws.com/bucket/json_examples//**/*.json` matches all the JSON files recursively.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-config
- Community: [Slack Channel](https://steampipe.io/community/join)
