# Table: config_ini_section

Retrieves all sections defined in an INI file.

## Examples

### Basic info

```sql
select
  path,
  section,
  comment
from
  config_ini_section
```

### List sections for a specific file

```sql
select
  path,
  section,
  comment
from
  config_ini_section
where
  path = '/Users/myuser/configs/main.ini';
```

### List subsections for a specific section

```sql
select
  path,
  section,
  comment
from
  config_ini_section
where
  section like 'settings.%';
```
