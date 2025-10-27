---
title: "Steampipe Table: xml_file - Query Config XML Files using SQL"
description: "Allows users to query XML Files in Config, specifically the file content in XML format, providing insights into configuration data and potential inconsistencies."
---

# Table: xml_file - Query Config XML Files using SQL

Config is a service that allows you to assess, audit, and evaluate the configurations of your config resources. It provides a centralized way to manage and evaluate configurations, including XML files, across your config resources. Config helps you maintain the desired state of your resources and take appropriate actions when predefined conditions are met.

## Table Usage Guide

The `xml_file` table provides insights into XML files within Config. As a DevOps engineer, explore file-specific details through this table, including content, file paths, and associated metadata. Utilize it to uncover information about XML files, such as their content and location, and to verify the consistency of configuration data.

For instance, if `xml_paths` is set to `[ "/Users/myuser/*.xml" ]`, and that directory contains:
- invoice.xml
- test.xml

This table will retrieve the file contents from each file mentioned above, which you can then query directly:

```sql
select
  path,
  jsonb_pretty(content) as file_content
from
  xml_file;
```

```sh
+---------------------------+----------------------------------------------------------------+
| path                      | file_content                                                   |
+---------------------------+----------------------------------------------------------------+
| /Users/myuser/test.xml    | {                                                              |
|                           |     "root": {                                                  |
|                           |         "foo": "bar",                                          |
|                           |         "includes": "common.xml"                               |
|                           |     }                                                          |
|                           | }                                                              |
| /Users/myuser/invoice.xml | {                                                              |
|                           |     "invoice": {                                               |
|                           |         "city": "East Centerville",                            |
|                           |         "date": "2012-08-06T00:00:00Z",                        |
|                           |         "items": [                                             |
|                           |             {                                                  |
|                           |                 "price": "1.47",                               |
|                           |                 "part_no": "A4786",                            |
|                           |                 "quantity": "4",                               |
|                           |                 "description": "Water Bucket (Filled)"         |
|                           |             },                                                 |
|                           |             {                                                  |
|                           |                 "size": "8",                                   |
|                           |                 "price": "133.7",                              |
|                           |                 "part_no": "E1628",                            |
|                           |                 "quantity": "1",                               |
|                           |                 "description": "High Heeled \"Ruby\" Slippers" |
|                           |             }                                                  |
|                           |         ],                                                     |
|                           |         "state": "KS",                                         |
|                           |         "street": "123 Tornado Alley\nSuite 16",               |
|                           |         "bill-to": "",                                         |
|                           |         "receipt": "Oz-Ware Purchase Invoice",                 |
|                           |         "ship-to": "",                                         |
|                           |         "customer": {                                          |
|                           |             "first_name": "Dorothy",                           |
|                           |             "family_name": "Gale"                              |
|                           |         }                                                      |
|                           |     }                                                          |
|                           | }                                                              |
+---------------------------+----------------------------------------------------------------+
```

or, you can query configurations of a particular file using:

```sql
select
  path,
  jsonb_pretty(content) as file_content
from
  xml_file
where
  path = '/Users/myuser/invoice.xml';
```

```sh
+----------------------------+----------------------------------------------------------------+
| path                       | file_content                                                   |
+----------------------------+----------------------------------------------------------------+
| ./Users/myuser/invoice.xml | {                                                              |
|                            |     "invoice": {                                               |
|                            |         "city": "East Centerville",                            |
|                            |         "date": "2012-08-06T00:00:00Z",                        |
|                            |         "items": [                                             |
|                            |             {                                                  |
|                            |                 "price": "1.47",                               |
|                            |                 "part_no": "A4786",                            |
|                            |                 "quantity": "4",                               |
|                            |                 "description": "Water Bucket (Filled)"         |
|                            |             },                                                 |
|                            |             {                                                  |
|                            |                 "size": "8",                                   |
|                            |                 "price": "133.7",                              |
|                            |                 "part_no": "E1628",                            |
|                            |                 "quantity": "1",                               |
|                            |                 "description": "High Heeled \"Ruby\" Slippers" |
|                            |             }                                                  |
|                            |         ],                                                     |
|                            |         "state": "KS",                                         |
|                            |         "street": "123 Tornado Alley\nSuite 16",               |
|                            |         "bill-to": "",                                         |
|                            |         "receipt": "Oz-Ware Purchase Invoice",                 |
|                            |         "ship-to": "",                                         |
|                            |         "customer": {                                          |
|                            |             "first_name": "Dorothy",                           |
|                            |             "family_name": "Gale"                              |
|                            |         }                                                      |
|                            |     }                                                          |
|                            | }                                                              |
+----------------------------+----------------------------------------------------------------+
```

## Examples

### Query a simple file
This query is useful for gaining insights into customer behavior by analyzing their purchase history. It enables you to identify the date of purchase, the customer's name, and the quantity of items ordered, which can be instrumental in understanding customer preferences and trends.
Given the file `invoice.xml` with the following configuration:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<invoice>
   <bill-to />
   <city>East Centerville</city>
   <customer>
      <family_name>Gale</family_name>
      <first_name>Dorothy</first_name>
   </customer>
   <date>2012-08-06T00:00:00Z</date>
   <items>
      <description>Water Bucket (Filled)</description>
      <part_no>A4786</part_no>
      <price>1.47</price>
      <quantity>4</quantity>
   </items>
   <items>
      <description>High Heeled "Ruby" Slippers</description>
      <part_no>E1628</part_no>
      <price>133.7</price>
      <quantity>1</quantity>
      <size>8</size>
   </items>
   <receipt>Oz-Ware Purchase Invoice</receipt>
   <ship-to />
   <state>KS</state>
   <street>123 Tornado Alley
Suite 16</street>
</invoice>
```

You can query the customer details and the number of items ordered:


```sql+postgres
select
  content -> 'invoice' ->> 'date' as order_date,
  concat(content -> 'invoice' -> 'customer' ->> 'first_name', ' ', content -> 'invoice' -> 'customer' ->> 'family_name') as customer_name,
  jsonb_array_length(content -> 'invoice' -> 'items') as order_count
from
  xml_file
where
  path = '/Users/myuser/invoice.xml';
```

```sql+sqlite
select
  json_extract(content, '$.invoice.date') as order_date,
  json_extract(content, '$.invoice.customer.first_name') || ' ' || json_extract(content, '$.invoice.customer.family_name') as customer_name,
  json_array_length(json_extract(content, '$.invoice.items')) as order_count
from
  xml_file
where
  path = '/Users/myuser/invoice.xml';
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
  (content -> 'invoice' ->> 'date')::timestamp as order_date,
  concat(content -> 'invoice' -> 'customer' ->> 'first_name', ' ', content -> 'invoice' -> 'customer' ->> 'family_name') as customer_name,
  item ->> 'description' as description,
  (item ->> 'price')::float as price,
  (item ->> 'quantity')::integer as quantity,
  (item ->> 'price')::float * (item ->> 'quantity')::integer as total
from
  xml_file,
  jsonb_array_elements(content -> 'invoice' -> 'items') as item
where
  path = '/Users/myuser/invoice.xml';
```

```sql+sqlite
select
  datetime(json_extract(content, '$.invoice.date')) as order_date,
  json_extract(content, '$.invoice.customer.first_name') || ' ' || json_extract(content, '$.invoice.customer.family_name') as customer_name,
  json_extract(item.value, '$.description') as description,
  cast(json_extract(item.value, '$.price') as real) as price,
  cast(json_extract(item.value, '$.quantity') as integer) as quantity,
  cast(json_extract(item.value, '$.price') as real) * cast(json_extract(item.value, '$.quantity') as integer) as total
from
  xml_file,
  json_each(json_extract(content, '$.invoice.items')) as item
where
  path = '/Users/myuser/invoice.xml';
```

```sh
+---------------------+---------------+-----------------------------+-------+----------+-------+
| order_date          | customer_name | description                 | price | quantity | total |
+---------------------+---------------+-----------------------------+-------+----------+-------+
| 2012-08-06 00:00:00 | Dorothy Gale  | Water Bucket (Filled)       | 1.47  | 4        | 5.88  |
| 2012-08-06 00:00:00 | Dorothy Gale  | High Heeled "Ruby" Slippers | 133.7 | 1        | 133.7 |
+---------------------+---------------+-----------------------------+-------+----------+-------+
```
