# Table: yml_file

Query the file contents from YML files found in the configured `yml_paths`.

For instance, if `yml_paths` is set to `[ "/Users/myuser/*.yml", "/Users/myuser/*.yaml" ]`, and that directory contains:

- invoice.yml
- test.yaml

This table will retrieve the file contents in JSON format from each file mentioned above, which you can then query directly:

```sql
select
  path,
  jsonb_pretty(content) as file_content
from
  yml_file;
```

```sh
+---------------------------+------------------------------------------------------------+
| path                      | file_content                                               |
+---------------------------+------------------------------------------------------------+
| /Users/myuser/invoice.yml | {                                                          |
|                           |     "city": "East Centerville",                            |
|                           |     "date": "2012-08-06T00:00:00Z",                        |
|                           |     "items": [                                             |
|                           |         {                                                  |
|                           |             "price": 1.47,                                 |
|                           |             "part_no": "A4786",                            |
|                           |             "quantity": 4,                                 |
|                           |             "description": "Water Bucket (Filled)"         |
|                           |         },                                                 |
|                           |         {                                                  |
|                           |             "size": 8,                                     |
|                           |             "price": 133.7,                                |
|                           |             "part_no": "E1628",                            |
|                           |             "quantity": 1,                                 |
|                           |             "description": "High Heeled \"Ruby\" Slippers" |
|                           |         }                                                  |
|                           |     ],                                                     |
|                           |     "state": "KS",                                         |
|                           |     "street": "123 Tornado Alley\nSuite 16\n",             |
|                           |     "bill-to": null,                                       |
|                           |     "receipt": "Oz-Ware Purchase Invoice",                 |
|                           |     "ship-to": null,                                       |
|                           |     "customer": {                                          |
|                           |         "first_name": "Dorothy",                           |
|                           |         "family_name": "Gale"                              |
|                           |     },                                                     |
|                           | }                                                          |
| /Users/myuser/test.yaml   | {                                                          |
|                           |     "foo": "bar",                                          |
|                           |     "includes": [                                          |
|                           |         "common.yaml"                                      |
|                           |     ]                                                      |
|                           | }                                                          |
+---------------------------+------------------------------------------------------------+
```

or, you can query configurations of a particular file using:

```sql
select
  path,
  jsonb_pretty(content) as file_content
from
  yml_file
where
  path = '/Users/myuser/invoice.yml';
```

```sh
+---------------------------+------------------------------------------------------------+
| path                      | file_content                                               |
+---------------------------+------------------------------------------------------------+
| /Users/myuser/invoice.yml | {                                                          |
|                           |     "city": "East Centerville",                            |
|                           |     "date": "2012-08-06T00:00:00Z",                        |
|                           |     "items": [                                             |
|                           |         {                                                  |
|                           |             "price": 1.47,                                 |
|                           |             "part_no": "A4786",                            |
|                           |             "quantity": 4,                                 |
|                           |             "description": "Water Bucket (Filled)"         |
|                           |         },                                                 |
|                           |         {                                                  |
|                           |             "size": 8,                                     |
|                           |             "price": 133.7,                                |
|                           |             "part_no": "E1628",                            |
|                           |             "quantity": 1,                                 |
|                           |             "description": "High Heeled \"Ruby\" Slippers" |
|                           |         }                                                  |
|                           |     ],                                                     |
|                           |     "state": "KS",                                         |
|                           |     "street": "123 Tornado Alley\nSuite 16\n",             |
|                           |     "bill-to": null,                                       |
|                           |     "receipt": "Oz-Ware Purchase Invoice",                 |
|                           |     "ship-to": null,                                       |
|                           |     "customer": {                                          |
|                           |         "first_name": "Dorothy",                           |
|                           |         "family_name": "Gale"                              |
|                           |     },                                                     |
|                           | }                                                          |
+---------------------------+------------------------------------------------------------+
```

## Examples

### Query a simple file

Given the file `invoice.yml` with the following configuration:

```yaml
---
receipt: Oz-Ware Purchase Invoice
date: 2012-08-06
customer:
  first_name: Dorothy
  family_name: Gale
items:
  - part_no: A4786
    description: Water Bucket (Filled)
    price: 1.47
    quantity: 4

  - part_no: E1628
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

You can query the customer details and the number of items ordered:

```sql
select
  content ->> 'date' as order_date,
  concat(content -> 'customer' ->> 'first_name', ' ', content -> 'customer' ->> 'family_name') as customer_name,
  jsonb_array_length(content -> 'items') as order_count
from
  yml_file
where
  path = '/Users/myuser/invoice.yml';
```

```sh
+----------------------+---------------+-------------+
| order_date           | customer_name | order_count |
+----------------------+---------------+-------------+
| 2012-08-06T00:00:00Z | Dorothy Gale  | 2           |
+----------------------+---------------+-------------+
```

### Casting column data for analysis

Text columns can be easily cast to other types:

```sql
select
  (content ->> 'date')::timestamp as order_date,
  concat(content -> 'customer' ->> 'first_name', ' ', content -> 'customer' ->> 'family_name') as customer_name,
  item ->> 'description' as description,
  (item ->> 'price')::float as price,
  (item ->> 'quantity')::integer as quantity,
  (item ->> 'price')::float * (item ->> 'quantity')::integer as total
from
  yml_file,
  jsonb_array_elements(content -> 'items') as item
where
  path = '/Users/myuser/invoice.yml';
```

```sh
+---------------------+---------------+-----------------------------+-------+----------+-------+
| order_date          | customer_name | description                 | price | quantity | total |
+---------------------+---------------+-----------------------------+-------+----------+-------+
| 2012-08-06 00:00:00 | Dorothy Gale  | Water Bucket (Filled)       | 1.47  | 4        | 5.88  |
| 2012-08-06 00:00:00 | Dorothy Gale  | High Heeled "Ruby" Slippers | 133.7 | 1        | 133.7 |
+---------------------+---------------+-----------------------------+-------+----------+-------+
```
