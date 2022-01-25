# Table: config_ini_key_value

Query data from INI files. The table will store all configurations along with sections of all INI file found in the configured `paths`.

For instance, if `paths` is set to `/Users/myuser/ini/*`, and that directory contains:

- defaults.ini
- sample.ini

This table will store all the configurations from each file mentioned above, along with section details, and comments.

Which you can then query directly:

```sql
select
  section,
  key,
  value
from
  config_ini_key_value;
```

```sh
+---------------------------------------------------------------+-------------------------------+-------------------------------+---------------------------+
| path                                                          | section                       | key                           | value                     |
+---------------------------------------------------------------+-------------------------------+-------------------------------+---------------------------+
| /Users/subhajit/Downloads/node_test/sample_files/sample.ini   | server                        | enforce_domain                | true                      |
| /Users/subhajit/Downloads/node_test/sample_files/defaults.ini | DEFAULT                       | instance_name                 | my-instance               |
| /Users/subhajit/Downloads/node_test/sample_files/defaults.ini | security                      | admin_user                    | admin                     |
| /Users/subhajit/Downloads/node_test/sample_files/sample.ini   | paths                         | data                          | /home/git/grafana         |
| /Users/subhajit/Downloads/node_test/sample_files/defaults.ini | plugin.grafana-image-renderer | rendering_ignore_https_errors | true                      |
| /Users/subhajit/Downloads/node_test/sample_files/defaults.ini | analytics                     | check_for_updates             | false                     |
| /Users/subhajit/Downloads/node_test/sample_files/sample.ini   | server                        | host                          | http://localhost:9999/api |
| /Users/subhajit/Downloads/node_test/sample_files/sample.ini   | server                        | http_port                     | 9999                      |
| /Users/subhajit/Downloads/node_test/sample_files/sample.ini   | profile testing               | aws_secret_access_key         | bar                       |
| /Users/subhajit/Downloads/node_test/sample_files/defaults.ini | auth.google                   | client_secret                 | 0ldS3cretKey              |
| /Users/subhajit/Downloads/node_test/sample_files/sample.ini   | server                        | protocol                      | http                      |
| /Users/subhajit/Downloads/node_test/sample_files/defaults.ini | database                      | port                          | 8080                      |
| /Users/subhajit/Downloads/node_test/sample_files/sample.ini   | profile testing               | aws_access_key_id             | foo                       |
| /Users/subhajit/Downloads/node_test/sample_files/sample.ini   | DEFAULT                       | app_mode                      | development               |
| /Users/subhajit/Downloads/node_test/sample_files/defaults.ini | database                      | url                           | http://localhost:8080/    |
+---------------------------------------------------------------+-------------------------------+-------------------------------+---------------------------+
```

or, you can query configurations of a particular file using:

```sql
select
  section,
  key,
  value
from
  config_ini_key_value
where
  path = '/Users/myuser/ini/defaults.ini';
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

and, the query is:

```sql
select
  section,
  key,
  value
from
  config_ini_key_value
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

or, you can check for value of particular keys:

```sql
select
  key,
  value
from
  config_ini_key_value
where
  path = '/Users/myuser/ini/defaults.ini'
  and section = 'analytics'
  and key = 'check_for_updates';
```

### Casting column data for analysis

Text columns can be easily cast to other types:

```sql
select
  key,
  value::bool
from
  config_ini_key_value
where
  path = '/Users/myuser/ini/defaults.ini'
  and section = 'analytics'
  and key = 'check_for_updates'
  and not value::bool;
```

### Query file with interpolation of values

Given the file `defaults.ini`with following configuration:

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

and, the environment variable `HOSTNAME` configured with:

```sh
export HOSTNAME=my-instance
```

In the above INI file, `instance_name` value refers to an environment variable `${HOSTNAME}`, and `url` refers to other values in the same section.
When querying the specific field, the table will store the actual value of `${HOSTNAME}`, i.e. `my-instance`.

```sql
select
  key,
  value
from
  config_ini_key_value
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

### Query files with nested value

Given the file `sample.ini`with following configuration:

```bash
[profile]
access_key = foo
secret_key = bar
s3 =
  max_concurrent_requests = 10
  max_queue_size = 1000
```

```sql
select
  section,
  key,
  value
from
  config_ini_key_value
where
  path = '/Users/myuser/ini/defaults.ini';
```

```sh
+---------+----------------------------+-------+
| section | key                        | value |
+---------+----------------------------+-------+
| profile | s3.max_concurrent_requests | 10    |
| profile | s3.max_queue_size          | 1000  |
| profile | access_key                 | foo   |
| profile | secret_key                 | bar   |
+---------+----------------------------+-------+
```
