![image](https://hub.steampipe.io/images/plugins/turbot/config-social-graphic.png)

# Config Plugin for Steampipe

Use SQL to query data from various types of configuration files, e.g. `INI`, `JSON`, `YML`.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/config)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/config/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-config/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install config
```

Configure your [config file](https://hub.steampipe.io/plugins/turbot/config#configuration) for each supported file type to include directories with files to be parsed. For any of the `paths` arguments, if no directory is specified, the current working directory will be used.

Run steampipe:

```shell
steampipe query
```

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

## Engines

This plugin is available for the following engines:

| Engine        | Description
|---------------|------------------------------------------
| [Steampipe](https://steampipe.io/docs) | The Steampipe CLI exposes APIs and services as a high-performance relational database, giving you the ability to write SQL-based queries to explore dynamic data. Mods extend Steampipe's capabilities with dashboards, reports, and controls built with simple HCL. The Steampipe CLI is a turnkey solution that includes its own Postgres database, plugin management, and mod support.
| [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview) | Steampipe Postgres FDWs are native Postgres Foreign Data Wrappers that translate APIs to foreign tables. Unlike Steampipe CLI, which ships with its own Postgres server instance, the Steampipe Postgres FDWs can be installed in any supported Postgres database version.
| [SQLite Extension](https://steampipe.io/docs//steampipe_sqlite/overview) | Steampipe SQLite Extensions provide SQLite virtual tables that translate your queries into API calls, transparently fetching information from your API or service as you request it.
| [Export](https://steampipe.io/docs/steampipe_export/overview) | Steampipe Plugin Exporters provide a flexible mechanism for exporting information from cloud services and APIs. Each exporter is a stand-alone binary that allows you to extract data using Steampipe plugins without a database.
| [Turbot Pipes](https://turbot.com/pipes/docs) | Turbot Pipes is the only intelligence, automation & security platform built specifically for DevOps. Pipes provide hosted Steampipe database instances, shared dashboards, snapshots, and more.

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-config.git
cd steampipe-plugin-config
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```sh
make
```

Configure the plugin:

```shell
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/config.spc
```

Try it!

```shell
steampipe query
> .inspect config
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Steampipe](https://steampipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #steampipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Config Plugin](https://github.com/turbot/steampipe-plugin-config/labels/help%20wanted)
