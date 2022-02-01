# Table: yml_key_value

Query key-value pair data along with comments and line number from YML files found in the configured `paths`.

For instance, if `paths` is set to `/Users/myuser/yml/*`, and that directory contains:

- sample.yml
- invoice.yml

This table will retrieve all key-value pairs from each file mentioned above, along with comments and the line number, which you can then query directly:

```sql
select
  key_path,
  value,
  tag as value_type,
  pre_comments,
  start_line
from
  yml_key_value;
```

```sh
+----------------------+---------------------------------------------------------------------------------------------------+------------+-----------------------------+------------+
| key_path             | value                                                                                             | value_type | pre_comments                | start_line |
+----------------------+---------------------------------------------------------------------------------------------------+------------+-----------------------------+------------+
| items.1.size         | 8                                                                                                 | integer    | []                          | 15         |
| specialDelivery      | Follow the Yellow Brick Road to the Emerald City. Pay no attention to the man behind the curtain. | string     | []                          | 25         |
| customer.family_name | Gale                                                                                              | string     | []                          | 6          |
| items.1.part_no      | E1628                                                                                             | string     | []                          | 13         |
| items.0.part_no      | A4786                                                                                             | string     | ["# List of ordered items"] | 9          |
| items.0.price        | 1.47                                                                                              | number     | []                          | 11         |
| date                 | 2012-08-06                                                                                        | string     | []                          | 3          |
| items.1.price        | 133.7                                                                                             | number     | []                          | 16         |
| customer.first_name  | Dorothy                                                                                           | string     | []                          | 5          |
| includes.0           | common.yaml                                                                                       | string     | []                          | 3          |
| foo                  | bar                                                                                               | string     | []                          | 4          |
| items.1.description  | High Heeled "Ruby" Slippers                                                                       | string     | []                          | 14         |
| receipt              | Oz-Ware Purchase Invoice                                                                          | string     | []                          | 2          |
| items.0.quantity     | 4                                                                                                 | integer    | []                          | 12         |
| items.0.description  | Water Bucket (Filled)                                                                             | string     | []                          | 10         |
| city                 | East Centerville                                                                                  | string     | []                          | 22         |
| items.1.quantity     | 1                                                                                                 | integer    | []                          | 17         |
| state                | KS                                                                                                | string     | []                          | 23         |
| bill_to              | <null>                                                                                            | null       | []                          | 18         |
| street               | 123 Tornado Alley                                                                                 | string     | []                          | 19         |
|                      | Suite 16                                                                                          |            |                             |            |
+----------------------+---------------------------------------------------------------------------------------------------+------------+-----------------------------+------------+
```

or, you can query configurations of a particular file using:

```sql
select
  key_path,
  value,
  tag as value_type,
  pre_comments,
  start_line
from
  yml_key_value
where
  path = '/Users/myuser/yml/invoice.yml';
```

## Examples

### Query a simple file

Given the file `invoice.yml` with following configuration:

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
specialDelivery: >
  Follow the Yellow Brick
  Road to the Emerald City.
  Pay no attention to the
  man behind the curtain.
```

and the query is:

```sql
with items as (
  select
    regexp_replace(key_path, '\.[a-z_]+', '') as item,
    regexp_replace(key_path, 'items\.[0-9]+\.', '') as data,
    value
  from
    yml_key_value
  where
    path = '/Users/myuser/yml/invoice.yml'
    and key_path like 'items.%'
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
+-----------------------------+---------+----------+-------+--------+
| item_name                   | part_no | quantity | price | size   |
+-----------------------------+---------+----------+-------+--------+
| Water Bucket (Filled)       | A4786   | 4        | 1.47  | <null> |
| High Heeled "Ruby" Slippers | E1628   | 1        | 133.7 | 8      |
+-----------------------------+---------+----------+-------+--------+
```

or, you can check the value for a particular key:

```sql
with items as (
  select
    regexp_replace(key_path, '\.[a-z_]+', '') as item,
    regexp_replace(key_path, 'items\.[0-9]+\.', '') as data,
    value
  from
    yml_key_value
  where
    path = '/Users/myuser/yml/invoice.yml'
    and key_path like 'items.%'
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

### Casting column data for analysis

Text columns can be easily cast to other types:

```sql
with items as (
  select
    regexp_replace(key_path, '\.[a-z_]+', '') as item,
    regexp_replace(key_path, 'items\.[0-9]+\.', '') as data,
    value
  from
    yml_key_value
  where
    path = '/Users/myuser/yml/invoice.yml'
    and key_path like 'items.%'
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
