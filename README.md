![image](https://hub.steampipe.io/images/plugins/turbot/config-social-graphic.png)

# Config Plugin for Steampipe

Use SQL to query data from various types of configuration files, e.g. `INI`, `JSON`, `YML`.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/config)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/config/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-config/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install config
```

Configure your [config file](https://hub.steampipe.io/plugins/turbot/config#configuration).

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

or, you can query configurations of a particular file using:

```sql
select
  path,
  section,
  key,
  value
from
  ini_key_value
where
  path = '/Users/myuser/ini/defaults.ini';
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

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-config/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Config Plugin](https://github.com/turbot/steampipe-plugin-config/labels/help%20wanted)
