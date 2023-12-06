---
title: "Steampipe Table: yml_key_value - Query Config YML Key Values using SQL"
description: "Allows users to query YML Key Values in Config, specifically the keys and their corresponding values in a YAML file, providing insights into configuration details and potential discrepancies."
---

# Table: yml_key_value - Query Config YML Key Values using SQL

Config is a service that enables you to assess, audit, and evaluate the configurations of your config resources. It provides a real-time snapshot of your resource configurations and lets you monitor configuration changes over time. Config helps you to ensure that your resources comply with your corporate standards and best practices.

## Table Usage Guide

The `yml_key_value` table provides insights into key-value pairs within YAML files in Config. As a DevOps engineer, explore key-specific details through this table, including their corresponding values and associated metadata. Utilize it to uncover information about keys, such as their hierarchal structure, the relationships between keys, and the verification of key values.

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
Explore a specific key-value pair in a YAML file to quickly identify a particular item's part number. This can be particularly useful for inventory management and tracking.
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
Determine the areas in your system where certain keys are less than a specified value. This practical application can be useful in scenarios where you need to filter out specific segments of your data based on your set criteria.
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
Explore specific parts within an invoice file to identify their unique part numbers. This is particularly useful when you need to quickly locate and assess individual parts within a large inventory.
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
Determine the details associated with a specific customer in an invoice file. This is useful for gaining insights into the customer's information, such as their first and last names.

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
This example demonstrates how to organize and examine data from a YAML file, specifically the details of different items from an invoice. The query allows for the easy examination of specific item details such as part number, item name, size, quantity, and price. Additionally, it provides a way to pinpoint information for a particular item using its part number, making it a valuable tool for inventory management and financial tracking.

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
This query is used to restructure and analyze invoice data stored in a YAML file. It allows users to understand the details of each item in the invoice, such as the part number, item name, size, quantity, and price, by transforming the data into a more readable and analyzable format.
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