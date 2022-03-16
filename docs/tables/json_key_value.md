# Table: json_key_value

Query key-value pair data along with comments and line number from JSON files found in the configured `paths`.

For instance, if `paths` is set to `/Users/myuser/json/*`, and that directory contains:

- sample.json
- invoice.json

This table will retrieve all key-value pairs from each file mentioned above, along with line numbers, which you can then query directly:

```sql
select
  key_path,
  value,
  start_line
from
  json_key_value;
```

```sh
+----------------------+-----------------------------+------------+
| key_path             | value                       | start_line |
+----------------------+-----------------------------+------------+
| items.1.part_no      | E1628                       | 19         |
| customer.first_name  | Dorothy                     | 6          |
| city                 | East Centerville            | 3          |
| items.1.size         | 8                           | 22         |
|                      |                             |            |
| items.1.price        | 133.7                       | 20         |
| items.1.quantity     | 1                           | 21         |
| street               | 123 Tornado Alley           | 28         |
|                      | Suite 16                    |            |
|                      |                             |            |
| state                | KS                          | 27         |
| items.0.description  | Water Bucket (Filled)       | 12         |
| items.0.price        | 1.47                        | 14         |
| date                 | 2012-08-06T00:00:00Z        | 9          |
| items.0.part_no      | A4786                       | 13         |
| items.1.description  | High Heeled "Ruby" Slippers | 18         |
| customer.family_name | Gale                        | 5          |
| receipt              | Oz-Ware Purchase Invoice    | 25         |
| items.0.quantity     | 4                           | 15         |
+----------------------+-----------------------------+------------+
```

or, you can query configurations of a particular file using:

```sql
select
  key_path,
  value,
  path
from
  json_key_value
where
  path = '/Users/myuser/json/invoice.json';
```

```sh
+----------------------+-----------------------------+---------------------------------+
| key_path             | value                       | path                            |
+----------------------+-----------------------------+---------------------------------+
| items.1.size         | 8                           | /Users/myuser/json/invoice.json |
| customer.family_name | Gale                        | /Users/myuser/json/invoice.json |
| items.1.part_no      | E1628                       | /Users/myuser/json/invoice.json |
| items.0.part_no      | A4786                       | /Users/myuser/json/invoice.json |
| items.0.price        | 1.47                        | /Users/myuser/json/invoice.json |
| date                 | 2012-08-06                  | /Users/myuser/json/invoice.json |
| items.1.price        | 133.7                       | /Users/myuser/json/invoice.json |
| customer.first_name  | Dorothy                     | /Users/myuser/json/invoice.json |
| includes.0           | common.yaml                 | /Users/myuser/json/invoice.json |
| foo                  | bar                         | /Users/myuser/json/invoice.json |
| items.1.description  | High Heeled "Ruby" Slippers | /Users/myuser/json/invoice.json |
| receipt              | Oz-Ware Purchase Invoice    | /Users/myuser/json/invoice.json |
| items.0.quantity     | 4                           | /Users/myuser/json/invoice.json |
| items.0.description  | Water Bucket (Filled)       | /Users/myuser/json/invoice.json |
| city                 | East Centerville            | /Users/myuser/json/invoice.json |
| items.1.quantity     | 1                           | /Users/myuser/json/invoice.json |
+----------------------+-----------------------------+---------------------------------+
```

## Examples

This table uses column `key_path` of type [ltree](https://www.postgresql.org/docs/12/ltree.html), contains a sequence of zero or more labels separated by dots representing a path from the root of a hierarchical tree to a particular node,  , which enables powerful search functionality that can be used to model, query and validate hierarchical and arbitrarily nested data structures, enables powerful search functionality that can be used to model, query and validate hierarchical and arbitrarily nested data structures.

### Search value of a particular key

```sql
select
  key_path,
  value as part_no
from
  json_key_value
where
  path = '/Users/myuser/json/invoice.json'
  and key_path = 'items.0.part_no';
```

or, you can query all subkeys that have index less than `items`,

```sql
select
  key_path,
  value as part_no
from
  json_key_value
where
  path = '/Users/myuser/json/invoice.json'
  and key_path < 'items';
```

```sh
+----------------------+----------------------+
| key_path             | value                |
+----------------------+----------------------+
| checks.0             | foo                  |
| checks.1             | timon                |
| city                 | East Centerville     |
| customer.family_name | Gale                 |
| customer.first_name  | Dorothy              |
| date                 | 2012-08-06T00:00:00Z |
+----------------------+----------------------+
```

### Search values using path matching

For example, from the sample JSON file above, you can query all `part_no` subkeys using the `~` operator to match an lquery,

```sql
select
  key_path,
  value as part_no
from
  json_key_value
where
  path = '/Users/myuser/json/invoice.json'
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

### Create a pivot table and search for specific key

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
