# Table: config_ini

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
  config_ini;
```

or, you can query configurations of a particular file using:

```sql
select
  section,
  key,
  value
from
  config_ini
where
  path = '/Users/myuser/ini/defaults.ini';
```

## Examples

### Query a simple file

Given the file `defaults.ini`with following configuration:

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
check_for_updates = false
```

and, the query is:

```sql
select
  section,
  key,
  value
from
  config_ini
where
  path = '/Users/myuser/ini/defaults.ini';
```

or, you can check for value of particular keys:

```sql
select
  key,
  value
from
  config_ini
where
  path = '/Users/myuser/ini/defaults.ini'
  and section = 'analytics'
  and key = 'check_for_updates';
```

### Query file with environment variable reference

Given the file `defaults.ini`with following configuration:

```bash
# default section
instance_name = ${HOSTNAME}

[security]
admin_user = admin

[auth.google]
client_secret = 0ldS3cretKey

[plugin.grafana-image-renderer]
rendering_ignore_https_errors = true
```

and, the environment variable `HOSTNAME` configured with:

```sh
export HOSTNAME=my-instance
```

In the above INI file, `instance_name` value refers to an environment variable `${HOSTNAME}`.
When querying the specific field, the table will store the actual value of `${HOSTNAME}`, i.e. `my-instance`.

```sql
select
  key,
  value
from
  config_ini
where
  path = '/Users/myuser/ini/defaults.ini';
```
