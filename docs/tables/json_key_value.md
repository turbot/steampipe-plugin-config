# Table: json_key_value

Query key-value pair data along with comments and line number from JSON files found in the configured `paths`.

For instance, if `paths` is set to `/Users/myuser/json/*`, and that directory contains:

- sample.json
- invoice.json

This table will retrieve all key-value pairs from each file mentioned above, along with the line number, which you can then query directly:

```sql
select
  key_path,
  value,
  tag as value_type,
  start_line
from
  json_key_value;
```

```sh
+----------------------+---------------------------------------------------------------------------------------------------+------------+------------+
| key_path             | value                                                                                             | value_type | start_line |
+----------------------+---------------------------------------------------------------------------------------------------+------------+------------+
| items.1.part_no      | E1628                                                                                             | string     | 19         |
| customer.first_name  | Dorothy                                                                                           | string     | 6          |
| city                 | East Centerville                                                                                  | string     | 3          |
| items.1.size         | 8                                                                                                 | integer    | 22         |
| specialDelivery      | Follow the Yellow Brick Road to the Emerald City. Pay no attention to the man behind the curtain. | string     | 26         |
|                      |                                                                                                   |            |            |
| items.1.price        | 133.7                                                                                             | number     | 20         |
| items.1.quantity     | 1                                                                                                 | integer    | 21         |
| street               | 123 Tornado Alley                                                                                 | string     | 28         |
|                      | Suite 16                                                                                          |            |            |
|                      |                                                                                                   |            |            |
| state                | KS                                                                                                | string     | 27         |
| items.0.description  | Water Bucket (Filled)                                                                             | string     | 12         |
| items.0.price        | 1.47                                                                                              | number     | 14         |
| date                 | 2012-08-06T00:00:00Z                                                                              | string     | 9          |
| items.0.part_no      | A4786                                                                                             | string     | 13         |
| items.1.description  | High Heeled "Ruby" Slippers                                                                       | string     | 18         |
| customer.family_name | Gale                                                                                              | string     | 5          |
| receipt              | Oz-Ware Purchase Invoice                                                                          | string     | 25         |
| items.0.quantity     | 4                                                                                                 | integer    | 15         |
+----------------------+---------------------------------------------------------------------------------------------------+------------+------------+
```

or, you can query configurations of a particular file using:

```sql
select
  key_path,
  value,
  tag as value_type,
  start_line
from
  json_key_value
where
  path = '/Users/myuser/json/invoice.json';
```

## Examples

### Query a simple file

Given the file `invoice.json` with following configuration:

```json
{
  "bill-to": null,
  "city": "East Centerville",
  "customer": { "family_name": "Gale", "first_name": "Dorothy" },
  "date": "2012-08-06T00:00:00Z",
  "items": [
    {
      "description": "Water Bucket (Filled)",
      "part_no": "A4786",
      "price": 1.47,
      "quantity": 4
    },
    {
      "description": "High Heeled \"Ruby\" Slippers",
      "part_no": "E1628",
      "price": 133.7,
      "quantity": 1,
      "size": 8
    }
  ],
  "receipt": "Oz-Ware Purchase Invoice",
  "ship-to": null,
  "specialDelivery": "Follow the Yellow Brick Road to the Emerald City. Pay no attention to the man behind the curtain.\n",
  "state": "KS",
  "street": "123 Tornado Alley\nSuite 16\n"
}
```

and the query is:

```sql
with items as (
  select
    subpath(key_path, 0, 2) as item,
    subpath(key_path, 2, 3) as data,
    value
  from
    json_key_value
  where
    path = '/Users/myuser/json/invoice.json'
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

or, you can check the value for a particular key:

```sql
with items as (
  select
    subpath(key_path, 0, 2) as item,
    subpath(key_path, 2, 3) as data,
    value
  from
    json_key_value
  where
    path = '/Users/myuser/json/invoice.json'
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

### Casting column data for analysis

Text columns can be easily cast to other types:

```sql
with items as (
  select
    subpath(key_path, 0, 2) as item,
    subpath(key_path, 2, 3) as data,
    value
  from
    json_key_value
  where
    path = '/Users/myuser/json/invoice.json'
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

### Searching key paths

The `Ltree` implements a materialized path, which is pretty quick for `SELECT` operations. This table uses `key_path` column of type `ltree` to represent the path in a hierarchical tree-like structure.

For example, from above mentioned sample JSON file,
you can easily query all the `part_no` using path match,

```sql
select
  value as part_no
from
  json_key_value
where
  path = '/Users/myuser/json/invoice.json'
  and key_path ~ 'items.*.part_no';
```

```sh
+---------+
| part_no |
+---------+
| E1628   |
| A4786   |
+---------+
```

or, you can fetch all details provided inside `customer` node,

```sql
select
  key_path,
  value
from
  json_key_value
where
  path = '/Users/myuser/json/invoice.json'
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
