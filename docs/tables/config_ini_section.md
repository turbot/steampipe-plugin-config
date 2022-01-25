# Table: config_ini_section

Retrieves all sections defined in an INI file.

## Examples

### Query sections of given file

```sql
select
  path,
  section,
  comment
from
  config_ini_section
where
  path = '/Users/myuser/ini/defaults.ini';
```

### List all configuration defined in default section

```sql
select
  sec.path as file_path,
  sec.section as section,
  f.key,
  f.value
from
  config_ini_section as sec,
  config_ini_key_value as f
where
  sec.path = '/Users/myuser/ini/defaults.ini'
  and f.path = sec.path
  and sec.section = 'DEFAULT';
```

```sh
+--------------------------------+---------+-------------------------------+------------------------+
| file_path                      | section | key                           | value                  |
+--------------------------------+---------+-------------------------------+------------------------+
| /Users/myuser/ini/defaults.ini | DEFAULT | url                           | http://localhost:8080/ |
| /Users/myuser/ini/defaults.ini | DEFAULT | check_for_updates             | false                  |
| /Users/myuser/ini/defaults.ini | DEFAULT | client_secret                 | 0ldS3cretKey           |
| /Users/myuser/ini/defaults.ini | DEFAULT | admin_user                    | admin                  |
| /Users/myuser/ini/defaults.ini | DEFAULT | instance_name                 | my-instance            |
| /Users/myuser/ini/defaults.ini | DEFAULT | port                          | 8080                   |
| /Users/myuser/ini/defaults.ini | DEFAULT | rendering_ignore_https_errors | true                   |
+--------------------------------+---------+-------------------------------+------------------------+
```
