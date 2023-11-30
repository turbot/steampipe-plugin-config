---
title: "Steampipe Table: yml_file - Query OCI Config YML Files using SQL"
description: "Allows users to query YML Files in OCI Config, specifically the details of each YML file, providing insights into file content and potential anomalies."
---

# Table: yml_file - Query OCI Config YML Files using SQL

Oracle Cloud Infrastructure (OCI) Config is a service that allows you to assess and evaluate the configurations of your OCI resources. It helps in managing, monitoring, and auditing resource configurations over time. An important part of this service is the YML Files, which contain the configuration details of various resources.

## Table Usage Guide

The `yml_file` table provides insights into YML Files within OCI Config. As a Cloud Architect or DevOps engineer, explore file-specific details through this table, including content, file size, and associated metadata. Utilize it to uncover information about files, such as those with specific configurations, the relationships between different configurations, and the verification of file content.

## Examples

### Query a simple file
This query is useful for gaining insights into customer purchase history from a specific file. It allows you to understand when and how many purchases a customer has made, which can help in analyzing customer behavior and trends.
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
Explore which items a customer has purchased and how much they've spent in total. This is useful for understanding individual customer behavior and identifying your best-selling products.
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