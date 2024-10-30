---
title: "Steampipe Table: toml_file - Query Config TOML Files using SQL"
description: "Allows users to query TOML Files in Config, specifically the file content in TOML format, providing insights into configuration data and potential inconsistencies."
---

# Table: toml_file - Query Config TOML Files using SQL

Config is a service that allows you to assess, audit, and evaluate the configurations of your config resources. It provides a centralized way to manage and evaluate configurations, including TOML files, across your config resources. Config helps you maintain the desired state of your resources and take appropriate actions when predefined conditions are met.

## Table Usage Guide

The `toml_file` table provides insights into TOML files within Config. As a DevOps engineer, explore file-specific details through this table, including content, file paths, and associated metadata. Utilize it to uncover information about TOML files, such as their content and location, and to verify the consistency of configuration data.

For instance, if `toml_paths` is set to `[ "/Users/myuser/*.toml" ]`, and that directory contains:
- invoice.toml
- test.toml

This table will retrieve the file contents from each file mentioned above, which you can then query directly:

```sql
select
  path,
  jsonb_pretty(content) as file_content
from
  toml_file;
```

```sh
+----------------------------+------------------------------------------------------------+
| path                       | file_content                                               |
+----------------------------+------------------------------------------------------------+
| /Users/myuser/test.toml    | {                                                          |
|                            |     "foo": "bar",                                          |
|                            |     "includes": [                                          |
|                            |         "common.toml"                                      |
|                            |     ]                                                      |
|                            | }                                                          |
| /Users/myuser/invoice.toml | {                                                          |
|                            |     "city": "East Centerville",                            |
|                            |     "date": "2012-08-06T00:00:00Z",                        |
|                            |     "items": [                                             |
|                            |         {                                                  |
|                            |             "price": 1.47,                                 |
|                            |             "part_no": "A4786",                            |
|                            |             "quantity": 4,                                 |
|                            |             "description": "Water Bucket (Filled)"         |
|                            |         },                                                 |
|                            |         {                                                  |
|                            |             "size": 8,                                     |
|                            |             "price": 133.7,                                |
|                            |             "part_no": "E1628",                            |
|                            |             "quantity": 1,                                 |
|                            |             "description": "High Heeled \"Ruby\" Slippers" |
|                            |         }                                                  |
|                            |     ],                                                     |
|                            |     "state": "KS",                                         |
|                            |     "street": "123 Tornado Alley\nSuite 16\n",             |
|                            |     "bill-to": "",                                         |
|                            |     "receipt": "Oz-Ware Purchase Invoice",                 |
|                            |     "ship-to": "",                                         |
|                            |     "customer": {                                          |
|                            |         "first_name": "Dorothy",                           |
|                            |         "family_name": "Gale"                              |
|                            |     }                                                      |
|                            | }                                                          |
+----------------------------+------------------------------------------------------------+
```

or, you can query configurations of a particular file using:

```sql
select
  path,
  jsonb_pretty(content) as file_content
from
  toml_file
where
  path = '/Users/myuser/invoice.toml';
```

```sh
+----------------------------+------------------------------------------------------------+
| path                       | file_content                                               |
+----------------------------+------------------------------------------------------------+
| /Users/myuser/invoice.toml | {                                                          |
|                            |     "city": "East Centerville",                            |
|                            |     "date": "2012-08-06T00:00:00Z",                        |
|                            |     "items": [                                             |
|                            |         {                                                  |
|                            |             "price": 1.47,                                 |
|                            |             "part_no": "A4786",                            |
|                            |             "quantity": 4,                                 |
|                            |             "description": "Water Bucket (Filled)"         |
|                            |         },                                                 |
|                            |         {                                                  |
|                            |             "size": 8,                                     |
|                            |             "price": 133.7,                                |
|                            |             "part_no": "E1628",                            |
|                            |             "quantity": 1,                                 |
|                            |             "description": "High Heeled \"Ruby\" Slippers" |
|                            |         }                                                  |
|                            |     ],                                                     |
|                            |     "state": "KS",                                         |
|                            |     "street": "123 Tornado Alley\nSuite 16\n",             |
|                            |     "bill-to": "",                                         |
|                            |     "receipt": "Oz-Ware Purchase Invoice",                 |
|                            |     "ship-to": "",                                         |
|                            |     "customer": {                                          |
|                            |         "first_name": "Dorothy",                           |
|                            |         "family_name": "Gale"                              |
|                            |     }                                                      |
|                            | }                                                          |
+----------------------------+------------------------------------------------------------+
```

## Examples

### Query a simple file
This query is useful for gaining insights into customer behavior by analyzing their purchase history. It enables you to identify the date of purchase, the customer's name, and the quantity of items ordered, which can be instrumental in understanding customer preferences and trends.
Given the file `invoice.toml` with the following configuration:

```toml
bill-to = ""
city = "East Centerville"

date = 2012-08-06T00:00:00Z

receipt = "Oz-Ware Purchase Invoice"
ship-to = ""
state = "KS"
street = """123 Tornado Alley
Suite 16
"""

[customer]
family_name = "Gale"
first_name = "Dorothy"

[[items]]
description = "Water Bucket (Filled)"
part_no = "A4786"
price = 1.47
quantity = 4

[[items]]
description = "High Heeled \"Ruby\" Slippers"
part_no = "E1628"
price = 133.7
quantity = 1
size = 8
```

You can query the customer details and the number of items ordered:


```sql+postgres
select
  content ->> 'date' as order_date,
  concat(content -> 'customer' ->> 'first_name', ' ', content -> 'customer' ->> 'family_name') as customer_name,
  jsonb_array_length(content -> 'items') as order_count
from
  toml_file
where
  path = '/Users/myuser/invoice.toml';
```

```sh
+----------------------+---------------+-------------+
| order_date           | customer_name | order_count |
+----------------------+---------------+-------------+
| 2012-08-06T00:00:00Z | Dorothy Gale  | 2           |
+----------------------+---------------+-------------+
```

### Casting column data for analysis
Determine the areas in which specific customer purchases can be analyzed. This allows for a detailed breakdown of individual transactions, including the customer's name, product descriptions, and the total cost of each item purchased.
Text columns can be easily cast to other types:


```sql+postgres
select
  (content ->> 'date')::timestamp as order_date,
  concat(content -> 'customer' ->> 'first_name', ' ', content -> 'customer' ->> 'family_name') as customer_name,
  item ->> 'description' as description,
  (item ->> 'price')::float as price,
  (item ->> 'quantity')::integer as quantity,
  (item ->> 'price')::float * (item ->> 'quantity')::integer as total
from
  toml_file,
  jsonb_array_elements(content -> 'items') as item
where
  path = '/Users/myuser/invoice.toml';
```

```sh
+---------------------+---------------+-----------------------------+-------+----------+-------+
| order_date          | customer_name | description                 | price | quantity | total |
+---------------------+---------------+-----------------------------+-------+----------+-------+
| 2012-08-06 00:00:00 | Dorothy Gale  | Water Bucket (Filled)       | 1.47  | 4        | 5.88  |
| 2012-08-06 00:00:00 | Dorothy Gale  | High Heeled "Ruby" Slippers | 133.7 | 1        | 133.7 |
+---------------------+---------------+-----------------------------+-------+----------+-------+
```
