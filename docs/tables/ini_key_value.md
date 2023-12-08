---
title: "Steampipe Table: ini_key_value - Query Config INI Key Values using SQL"
description: "Allows users to query INI Key Values in Config, specifically the key-value pairs in INI configuration files, enabling detailed analysis of configuration settings."
---

# Table: ini_key_value - Query Config INI Key Values using SQL

INI Key Values in Config are key-value pairs in INI configuration files. These are used to store configuration settings in a structured format that is easy to read and write. These key-value pairs can be used to configure the behavior of software applications, libraries, and drivers.

## Table Usage Guide

The `ini_key_value` table provides insights into key-value pairs in INI configuration files within Config. As a DevOps engineer, explore key-value pair details through this table, including keys, values, and associated metadata. Utilize it to uncover information about configuration settings, such as those for software applications, libraries, and drivers.

## Examples

### Query a simple file
This query allows you to explore the contents of a specific configuration file, which can help you understand and manage your application's settings. For instance, you might use this query to determine if certain security or analytics settings are configured as expected.
Given the file `defaults.ini` with following configuration:

```bash
# default section
instance_name = my-instance

[security]
admin_user = admin

[auth.google]
client_secret = 0ldS3cretKey

[plugin.grafana]
ignore_https_errors = true

[analytics]
check_for_updates = false
```

and the query is:


```sql+postgres
select
  section,
  key,
  value
from
  ini_key_value
where
  path = '/Users/myuser/defaults.ini';
```

```sql+sqlite
select
  section,
  key,
  value
from
  ini_key_value
where
  path = '/Users/myuser/defaults.ini';
```

```sh
+----------------+---------------------+--------------+
| section        | key                 | value        |
+----------------+---------------------+--------------+
| security       | admin_user          | admin        |
| DEFAULT        | instance_name       | my-instance  |
| auth.google    | client_secret       | 0ldS3cretKey |
| plugin.grafana | ignore_https_errors | true         |
| analytics      | check_for_updates   | false        |
+----------------+---------------------+--------------+
```

or, you can check the value for a particular key:

```sql+postgres
select
  section,
  key,
  value
from
  ini_key_value
where
  path = '/Users/myuser/defaults.ini'
  and section = 'analytics'
  and key = 'check_for_updates';
```

```sql+sqlite
select
  section,
  key,
  value
from
  ini_key_value
where
  path = '/Users/myuser/defaults.ini'
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
Determine the areas in which automatic updates for analytics are disabled. This is useful to ensure that all analytics tools are up-to-date and running the latest versions.
Text columns can be easily cast to other types:


```sql+postgres
select
  section,
  key,
  value::bool
from
  ini_key_value
where
  path = '/Users/myuser/defaults.ini'
  and section = 'analytics'
  and key = 'check_for_updates'
  and not value::bool;
```

```sql+sqlite
select
  section,
  key,
  value in ('t', 'true', '1')
from
  ini_key_value
where
  path = '/Users/myuser/defaults.ini'
  and section = 'analytics'
  and key = 'check_for_updates'
  and not value in ('t', 'true', '1');
```

```sh
+-----------+-------------------+-------+
| section   | key               | value |
+-----------+-------------------+-------+
| analytics | check_for_updates | false |
+-----------+-------------------+-------+
```

### Query a file with value interpolation
Determine the specific settings and configurations within a file to gain insights into its structure and content. This can be useful for understanding how a system or application is configured and identifying any potential issues or areas for improvement.
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

[plugin.grafana]
ignore_https_errors = true
```

and the environment variable `HOSTNAME` is set:

```sh
export HOSTNAME=my-instance
```

In the above INI file, the value for `instance_name` refers to an environment variable `${HOSTNAME}` and `url` refers to other another value in the same section.

When querying values, the table will store the interpolated values, e.g., `${HOSTNAME}` will be stored as `my-instance`.


```sql+postgres
select
  key,
  value
from
  ini_key_value
where
  path = '/Users/myuser/defaults.ini';
```

```sql+sqlite
select
  key,
  value
from
  ini_key_value
where
  path = '/Users/myuser/defaults.ini';
```

```sh
+----------------+---------------------+------------------------+
| section        | key                 | value                  |
+----------------+---------------------+------------------------+
| plugin.grafana | ignore_https_errors | true                   |
| database       | port                | 8080                   |
| DEFAULT        | instance_name       | my-instance            |
| analytics      | check_for_updates   | false                  |
| auth.google    | client_secret       | 0ldS3cretKey           |
| security       | admin_user          | admin                  |
| database       | url                 | http://localhost:8080/ |
+----------------+---------------------+------------------------+
```