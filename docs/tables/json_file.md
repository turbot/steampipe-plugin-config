# Table: json_file

Query the file contents from JSON files found in the configured `paths`.

For instance, if `paths` is set to `/Users/myuser/json/*`, and that directory contains:

- sample.json
- invoice.json

This table will retrieve the file contents from each file mentioned above, which you can then query directly:

```sql
select
  path,
  jsonb_pretty(content) as file_content
from
  json_file;
```

```sh
+---------------------------------+------------------------------------------------------------------------------------------------------------------------------+
| path                            | file_content                                                                                                                 |
+---------------------------------+------------------------------------------------------------------------------------------------------------------------------+
| /Users/myuser/json/invoice.json | {                                                                                                                            |
|                                 |     "city": "East Centerville",                                                                                              |
|                                 |     "date": "2012-08-06T00:00:00Z",                                                                                          |
|                                 |     "items": [                                                                                                               |
|                                 |         {                                                                                                                    |
|                                 |             "price": 1.47,                                                                                                   |
|                                 |             "part_no": "A4786",                                                                                              |
|                                 |             "quantity": 4,                                                                                                   |
|                                 |             "description": "Water Bucket (Filled)"                                                                           |
|                                 |         },                                                                                                                   |
|                                 |         {                                                                                                                    |
|                                 |             "size": 8,                                                                                                       |
|                                 |             "price": 133.7,                                                                                                  |
|                                 |             "part_no": "E1628",                                                                                              |
|                                 |             "quantity": 1,                                                                                                   |
|                                 |             "description": "High Heeled \"Ruby\" Slippers"                                                                   |
|                                 |         }                                                                                                                    |
|                                 |     ],                                                                                                                       |
|                                 |     "state": "KS",                                                                                                           |
|                                 |     "street": "123 Tornado Alley\nSuite 16\n",                                                                               |
|                                 |     "bill-to": null,                                                                                                         |
|                                 |     "receipt": "Oz-Ware Purchase Invoice",                                                                                   |
|                                 |     "ship-to": null,                                                                                                         |
|                                 |     "customer": {                                                                                                            |
|                                 |         "first_name": "Dorothy",                                                                                             |
|                                 |         "family_name": "Gale"                                                                                                |
|                                 |     },                                                                                                                       |
|                                 |     "specialDelivery": "Follow the Yellow Brick Road to the Emerald City. Pay no attention to the man behind the curtain.\n" |
|                                 | }                                                                                                                            |
| /Users/myuser/json/test.json    | {                                                                                                                            |
|                                 |     "foo": "bar",                                                                                                            |
|                                 |     "includes": [                                                                                                            |
|                                 |         "common.json"                                                                                                        |
|                                 |     ]                                                                                                                        |
|                                 | }                                                                                                                            |
+---------------------------------+------------------------------------------------------------------------------------------------------------------------------+
```

or, you can query configurations of a particular file using:

```sql
select
  path,
  jsonb_pretty(content) as file_content
from
  json_file
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

and you can easily query the customer details and the number of items ordered:

```sql
select
  content ->> 'date' as order_date,
  concat(content -> 'customer' ->> 'first_name', ' ', content -> 'customer' ->> 'family_name') as customer_name,
  jsonb_array_length(content -> 'items') as order_count
from
  json_file
where
  path = '/Users/myuser/json/invoice.json';
```

```sh
+----------------------+---------------+-------------+
| order_date           | customer_name | order_count |
+----------------------+---------------+-------------+
| 2012-08-06T00:00:00Z | Dorothy Gale  | 2           |
+----------------------+---------------+-------------+
```

or, you can also list the ordered items:

```sql
select
  content ->> 'date' as order_date,
  concat(content -> 'customer' ->> 'first_name', ' ', content -> 'customer' ->> 'family_name') as customer_name,
  item ->> 'description' as description,
  (item ->> 'price')::float as price,
  (item ->> 'quantity')::integer as quantity,
  (item ->> 'price')::float * (item ->> 'quantity')::integer as total
from
  json_file,
  jsonb_array_elements(content -> 'items') as item
where
  path = '/Users/myuser/json/invoice.json';
```

```sh
+----------------------+---------------+-----------------------------+-------+----------+-------+
| order_date           | customer_name | description                 | price | quantity | total |
+----------------------+---------------+-----------------------------+-------+----------+-------+
| 2012-08-06T00:00:00Z | Dorothy Gale  | Water Bucket (Filled)       | 1.47  | 4        | 5.88  |
| 2012-08-06T00:00:00Z | Dorothy Gale  | High Heeled "Ruby" Slippers | 133.7 | 1        | 133.7 |
+----------------------+---------------+-----------------------------+-------+----------+-------+
```

### Casting column data for analysis

Text columns can be easily cast to other types:

```sql
select
  content ->> 'date' as order_date,
  concat(content -> 'customer' ->> 'first_name', ' ', content -> 'customer' ->> 'family_name') as customer_name,
  item ->> 'description' as description,
  (item ->> 'price')::float as price,
  (item ->> 'quantity')::integer as quantity,
  (item ->> 'price')::float * (item ->> 'quantity')::integer as total
from
  json_file,
  jsonb_array_elements(content -> 'items') as item
where
  path = '/Users/myuser/json/invoice.json';
```
