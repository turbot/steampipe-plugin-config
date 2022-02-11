# Table: ini_key_value

Query section and key-value pair data from INI files found in the configured `ini_paths`.

For instance, if `ini_paths` is set to `[ "/Users/myuser/ini/*.ini" ]`, and that directory contains:

- defaults.ini
- sample.ini

This table will retrieve all key-value pairs from each file mentioned above, along with section names and comments, which you can then query directly:

```sql
select
  path,
  section,
  key,
  value,
  comment
from
  ini_key_value;
```

```sh
+--------------------------------+-------------------------------+-------------------------------+---------------------------+-----------------------------+
| path                           | section                       | key                           | value                     | comment                     |
+--------------------------------+-------------------------------+-------------------------------+---------------------------+-----------------------------+
| /Users/myuser/ini/defaults.ini | analytics                     | check_for_updates             | false                     |                             |
| /Users/myuser/ini/defaults.ini | auth.google                   | client_secret                 | 0ldS3cretKey              |                             |
| /Users/myuser/ini/defaults.ini | database                      | port                          | 8080                      |                             |
| /Users/myuser/ini/defaults.ini | database                      | url                           | http://localhost:8080/    |                             |
| /Users/myuser/ini/defaults.ini | DEFAULT                       | instance_name                 | my-instance               |                             |
| /Users/myuser/ini/defaults.ini | plugin.grafana-image-renderer | rendering_ignore_https_errors | true                      |                             |
| /Users/myuser/ini/defaults.ini | security                      | admin_user                    | admin                     |                             |
| /Users/myuser/ini/sample.ini   | DEFAULT                       | app_mode                      | development               |                             |
| /Users/myuser/ini/sample.ini   | paths                         | data                          | /home/git/grafana         |                             |
| /Users/myuser/ini/sample.ini   | profile testing               | aws_access_key_id             | foo                       |                             |
| /Users/myuser/ini/sample.ini   | profile testing               | aws_secret_access_key         | bar                       |                             |
| /Users/myuser/ini/sample.ini   | server                        | enforce_domain                | true                      |                             |
| /Users/myuser/ini/sample.ini   | server                        | host                          | http://localhost:9999/api | # Update after moving hosts |
| /Users/myuser/ini/sample.ini   | server                        | http_port                     | 9999                      |                             |
| /Users/myuser/ini/sample.ini   | server                        | protocol                      | http                      |                             |
+--------------------------------+-------------------------------+-------------------------------+---------------------------+-----------------------------+
```

or, you can query configurations of a particular file using:

```sql
select
  section,
  key,
  value,
  comment
from
  ini_key_value
where
  path = '/Users/myuser/ini/defaults.ini';
```

```sh
+--------------------------------+-------------------------------+-------------------------------+---------------------------+---------+
| path                           | section                       | key                           | value                     | comment |
+--------------------------------+-------------------------------+-------------------------------+---------------------------+---------+
| /Users/myuser/ini/defaults.ini | analytics                     | check_for_updates             | false                     |         |
| /Users/myuser/ini/defaults.ini | auth.google                   | client_secret                 | 0ldS3cretKey              |         |
| /Users/myuser/ini/defaults.ini | database                      | port                          | 8080                      |         |
| /Users/myuser/ini/defaults.ini | database                      | url                           | http://localhost:8080/    |         |
| /Users/myuser/ini/defaults.ini | DEFAULT                       | instance_name                 | my-instance               |         |
| /Users/myuser/ini/defaults.ini | plugin.grafana-image-renderer | rendering_ignore_https_errors | true                      |         |
| /Users/myuser/ini/defaults.ini | security                      | admin_user                    | admin                     |         |
+--------------------------------+-------------------------------+-------------------------------+---------------------------+---------+
```

## Examples

### Query a simple file

Given the file `defaults.ini` with following configuration:

```bash
# default section
instance_name = my-instance

[security]
admin_user = admin

[auth.google]
client_secret = 0ldS3cretKey

[plugin.grafana-image-renderer]
rendering_ignore_https_errors = true

[analytics]
# Set to false to disable all checks for new versions of installed plugins
check_for_updates = false
```

and the query is:

```sql
select
  section,
  key,
  value
from
  ini_key_value
where
  path = '/Users/myuser/ini/defaults.ini';
```

```sh
+-------------------------------+-------------------------------+--------------+
| section                       | key                           | value        |
+-------------------------------+-------------------------------+--------------+
| auth.google                   | client_secret                 | 0ldS3cretKey |
| analytics                     | check_for_updates             | false        |
| security                      | admin_user                    | admin        |
| plugin.grafana-image-renderer | rendering_ignore_https_errors | true         |
| DEFAULT                       | instance_name                 | my-instance  |
+-------------------------------+-------------------------------+--------------+
```

or, you can check the value for a particular key:

```sql
select,
  section,
  key,
  value
from
  ini_key_value
where
  path = '/Users/myuser/ini/defaults.ini'
  and section = 'analytics'
  and key = 'check_for_updates';
```

```sh
+-----------+-------------------+-------+
| section   | key               | value |
+-----------+-------------------+-------+
| analytics | check_for_updates | false |
+-----------+-------------------+-------+
```

### Casting column data for analysis

Text columns can be easily cast to other types:

```sql
select
  section,
  key,
  value::bool
from
  ini_key_value
where
  path = '/Users/myuser/ini/defaults.ini'
  and section = 'analytics'
  and key = 'check_for_updates'
  and not value::bool;
```

```sh
+-----------+-------------------+-------+
| section   | key               | value |
+-----------+-------------------+-------+
| analytics | check_for_updates | false |
+-----------+-------------------+-------+
```

### Query a file with value interpolation

Given the file `defaults.ini` with following configuration:

```bash
# default section
instance_name = ${HOSTNAME}

[security]
admin_user = admin

[auth.google]
client_secret = 0ldS3cretKey

[database]
port = 8080
url = http://localhost:%(port)s/

[plugin.grafana-image-renderer]
rendering_ignore_https_errors = true
```

and the environment variable `HOSTNAME`:

```sh
export HOSTNAME=my-instance
```

In the above INI file, the value for `instance_name` refers to an environment variable `${HOSTNAME}` and `url` refers to other another value in the same section.

When querying values, the table will store the interpolated values, e.g., `${HOSTNAME}` will be stored as `my-instance`.

```sql
select
  key,
  value
from
  ini_key_value
where
  path = '/Users/myuser/ini/defaults.ini';
```

```sh
+-------------------------------+-------------------------------+------------------------+
| section                       | key                           | value                  |
+-------------------------------+-------------------------------+------------------------+
| plugin.grafana-image-renderer | rendering_ignore_https_errors | true                   |
| database                      | port                          | 8080                   |
| DEFAULT                       | instance_name                 | my-instance            |
| analytics                     | check_for_updates             | false                  |
| auth.google                   | client_secret                 | 0ldS3cretKey           |
| security                      | admin_user                    | admin                  |
| database                      | url                           | http://localhost:8080/ |
+-------------------------------+-------------------------------+------------------------+
```
