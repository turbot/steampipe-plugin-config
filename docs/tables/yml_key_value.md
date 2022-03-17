# Table: yml_key_value

Query key-value pairs, comments, and line numbers from YML files found in the configured `yml_paths`.

For instance, if `yml_paths` is set to `[ "/Users/myuser/*.yml", "/Users/myuser/*.yaml" ]`, and that directory contains:

- sample.yml
- invoice.yml

This table will retrieve all key-value pairs from each file mentioned above, along with comments and line numbers, which you can then query directly:

```sql
select
  key_path,
  value,
  tag,
  pre_comments,
  start_line
from
  yml_key_value;
```

```sh
+----------------------+-----------------------------+---------+-----------------------------+------------+
| key_path             | value                       | tag     | pre_comments                | start_line |
+----------------------+-----------------------------+---------+-----------------------------+------------+
| items.1.size         | 8                           | !!int   | []                          | 15         |
| customer.family_name | Gale                        | !!str   | []                          | 6          |
| items.1.part_no      | E1628                       | !!str   | []                          | 13         |
| items.0.part_no      | A4786                       | !!str   | ["# List of ordered items"] | 9          |
| items.0.price        | 1.47                        | !!float | []                          | 11         |
| date                 | 2012-08-06                  | !!str   | []                          | 3          |
| items.1.price        | 133.7                       | !!float | []                          | 16         |
| customer.first_name  | Dorothy                     | !!str   | []                          | 5          |
| includes.0           | common.yaml                 | !!str   | []                          | 3          |
| foo                  | bar                         | !!str   | []                          | 4          |
| items.1.description  | High Heeled "Ruby" Slippers | !!str   | []                          | 14         |
| receipt              | Oz-Ware Purchase Invoice    | !!str   | []                          | 2          |
| items.0.quantity     | 4                           | !!int   | []                          | 12         |
| items.0.description  | Water Bucket (Filled)       | !!str   | []                          | 10         |
| city                 | East Centerville            | !!str   | []                          | 22         |
| items.1.quantity     | 1                           | !!int   | []                          | 17         |
| state                | KS                          | !!str   | []                          | 23         |
| bill_to              | <null>                      | !!null  | []                          | 18         |
| street               | 123 Tornado Alley           | !!str   | []                          | 19         |
|                      | Suite 16                    |         |                             |            |
+----------------------+-----------------------------+---------+-----------------------------+------------+
```

or, you can query configurations of a particular file using:

```sql
select
  key_path,
  value,
  path
from
  yml_key_value
where
  path = '/Users/myuser/yml/invoice.yml';
```

```sh
+----------------------+-----------------------------+-------------------------------+
| key_path             | value                       | path                          |
+----------------------+-----------------------------+-------------------------------+
| items.1.size         | 8                           | /Users/myuser/yml/invoice.yml |
| customer.family_name | Gale                        | /Users/myuser/yml/invoice.yml |
| items.1.part_no      | E1628                       | /Users/myuser/yml/invoice.yml |
| items.0.part_no      | A4786                       | /Users/myuser/yml/invoice.yml |
| items.0.price        | 1.47                        | /Users/myuser/yml/invoice.yml |
| date                 | 2012-08-06                  | /Users/myuser/yml/invoice.yml |
| items.1.price        | 133.7                       | /Users/myuser/yml/invoice.yml |
| customer.first_name  | Dorothy                     | /Users/myuser/yml/invoice.yml |
| includes.0           | common.yaml                 | /Users/myuser/yml/invoice.yml |
| foo                  | bar                         | /Users/myuser/yml/invoice.yml |
| items.1.description  | High Heeled "Ruby" Slippers | /Users/myuser/yml/invoice.yml |
| receipt              | Oz-Ware Purchase Invoice    | /Users/myuser/yml/invoice.yml |
| items.0.quantity     | 4                           | /Users/myuser/yml/invoice.yml |
| items.0.description  | Water Bucket (Filled)       | /Users/myuser/yml/invoice.yml |
| city                 | East Centerville            | /Users/myuser/yml/invoice.yml |
| items.1.quantity     | 1                           | /Users/myuser/yml/invoice.yml |
+----------------------+-----------------------------+-------------------------------+
```

## Examples

The `key_path` column's data type is
[ltree](https://www.postgresql.org/docs/12/ltree.html), so all `key_path`
values are stored as dot-delimited label paths. This enables the use of the
usual comparison operators along with `ltree` operators and functions which can
be used to match subpaths, find ancestors and descendants, and search arrays.

For all examples below, assume we're using the file `invoice.yml` with the following configuration:

```yml
---
receipt: Oz-Ware Purchase Invoice
date: 2012-08-06
customer:
    first_name:  Dorothy
    family_name: Gale
# List of ordered items
items:
  - part_no: A4786 # item 1
    description: Water Bucket (Filled)
    price: 1.47
    quantity: 4
  - part_no: E1628 # item 2
    description: High Heeled "Ruby" Slippers
    size: 8
    price: 133.7
    quantity: 1
bill-to: &id001
street: |
  123 Tornado Alley
  Suite 16
city: East Centerville
state: KS
ship-to: *id001
```

### Query a specific key-value pair

You can query a specific key path to get its value:

```sql
select
  key_path,
  value as part_no
from
  json_key_value
where
  path = '/Users/myuser/yml/invoice.yml'
  and key_path = 'items.0.part_no';
```

```sh
+-----------------+-------+
| key_path        | value |
+-----------------+-------+
| items.0.part_no | A4786 |
+-----------------+-------+
```

### Query using comparison operators

The usual comparison operators, like `<`, `>`, `<=`, and `>=` work with `ltree` columns.

For instance, you can use the `<` operator to query all key paths that are before `items` alphabetically:

```sql
select
  key_path,
  value as part_no
from
  json_key_value
where
  path = '/Users/myuser/yml/invoice.yml'
  and key_path < 'items';
```

```sh
+----------------------+----------------------+
| key_path             | part_no              |
+----------------------+----------------------+
| bill_to              | <null>               |
| city                 | East Centerville     |
| customer.family_name | Gale                 |
| customer.first_name  | Dorothy              |
| date                 | 2012-08-06T00:00:00Z |
+----------------------+----------------------+
```

### Query using path matching

`ltree` also supports additional operators like `~` which can be used to find all `part_no` subkeys:

```sql
select
  key_path,
  value as part_no
from
  json_key_value
where
  path = '/Users/myuser/yml/invoice.yml'
  and key_path ~ 'items.*.part_no';
```

```sh
+-----------------+---------+
| key_path        | part_no |
+-----------------+---------+
| items.1.part_no | E1628   |
| items.0.part_no | A4786   |
+-----------------+---------+
```

### List descendants of a specific node

```sql
select
  key_path,
  value
from
  json_key_value
where
  path = '/Users/myuser/yml/invoice.yml'
  and key_path <@ 'customer';
```

```sh
+----------------------+---------+
| key_path             | value   |
+----------------------+---------+
| customer.first_name  | Dorothy |
| customer.family_name | Gale    |
+----------------------+---------+
```

### Create a pivot table and search for a specific key

```sql
with items as (
  select
    subpath(key_path, 0, 2) as item,
    subpath(key_path, 2, 3) as data,
    value
  from
    yml_key_value
  where
    path = '/Users/myuser/yml/invoice.yml'
    and key_path ~ 'items.*'
)
select
  max(case when data = 'part_no' then value else null end) as part_no,
  max(case when data = 'description' then value else null end) as item_name,
  max(case when data = 'size' then value else null end) as size,
  max(case when data = 'quantity' then value else null end) as quantity,
  max(case when data = 'price' then value else null end) as price
from
  items
group by
  item;
```

```sh
+---------+-----------------------------+--------+----------+-------+
| part_no | item_name                   | size   | quantity | price |
+---------+-----------------------------+--------+----------+-------+
| A4786   | Water Bucket (Filled)       | <null> | 4        | 1.47  |
| E1628   | High Heeled "Ruby" Slippers | 8      | 1        | 133.7 |
+---------+-----------------------------+--------+----------+-------+
```

You can also check the value for a particular key:

```sql
with items as (
  select
    subpath(key_path, 0, 2) as item,
    subpath(key_path, 2, 3) as data,
    value
  from
    yml_key_value
  where
    path = '/Users/myuser/yml/invoice.yml'
    and key_path ~ 'items.*'
),
pivot_tables as (
  select
  max(case when data = 'part_no' then value else null end) as part_no,
  max(case when data = 'description' then value else null end) as item_name,
  max(case when data = 'size' then value else null end) as size,
  max(case when data = 'quantity' then value else null end) as quantity,
  max(case when data = 'price' then value else null end) as price
from
  items
group by
  item
)
select * from pivot_tables where part_no = 'E1628';
```

```sh
+---------+-----------------------------+--------+----------+-------+
| part_no | item_name                   | size   | quantity | price |
+---------+-----------------------------+--------+----------+-------+
| E1628   | High Heeled "Ruby" Slippers | 8      | 1        | 133.7 |
+---------+-----------------------------+--------+----------+-------+
```

### Casting column data for analysis

The `value` column data type is `text`, so you can easily cast it when required:

```sql
with items as (
  select
    subpath(key_path, 0, 2) as item,
    subpath(key_path, 2, 3) as data,
    value
  from
    yml_key_value
  where
    path = '/Users/myuser/yml/invoice.yml'
    and key_path ~ 'items.*'
)
select
  max(case when data = 'part_no' then value else null end) as part_no,
  max(case when data = 'description' then value else null end) as item_name,
  (max(case when data = 'size' then value else null end))::integer as size,
  (max(case when data = 'quantity' then value else null end))::integer as quantity,
  (max(case when data = 'price' then value else null end))::float as price
from
  items
group by
  item;
```

```sh
+---------+-----------------------------+--------+----------+-------+
| part_no | item_name                   | size   | quantity | price |
+---------+-----------------------------+--------+----------+-------+
| A4786   | Water Bucket (Filled)       | <null> | 4        | 1.47  |
| E1628   | High Heeled "Ruby" Slippers | 8      | 1        | 133.7 |
+---------+-----------------------------+--------+----------+-------+
```
