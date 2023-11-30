---
title: "Steampipe Table: ini_section - Query Config INI Sections using SQL"
description: "Allows users to query INI Sections in Config, specifically the section names and associated properties, providing insights into configuration settings and potential discrepancies."
---

# Table: ini_section - Query Config INI Sections using SQL

An INI Section is a part of the INI file format that is used for configuration settings. INI files are simple text files with a basic structure composed of sections, properties, and values. The sections, denoted by the square bracket syntax, categorize the properties and values for easy retrieval and readability.

## Table Usage Guide

The `ini_section` table provides insights into INI sections within configuration files. As a DevOps engineer, explore section-specific details through this table, including associated properties and values. Utilize it to uncover information about configuration settings, such as those with undefined properties, the organization of sections, and the verification of property values.

## Examples

### Basic info
Explore the structure and comments of your configuration files to better understand their organization and purpose. This can be particularly useful when troubleshooting or optimizing your system setup.

```sql
select
  path,
  section,
  comment
from
  ini_section;
```

### List sections for a specific file
Explore which sections and associated comments are present in a specific configuration file. This is useful for understanding the configuration and settings of your applications.

```sql
select
  path,
  section,
  comment
from
  ini_section
where
  path = '/Users/myuser/configs/main.ini';
```

### List subsections for a specific section
Discover the segments that fall under a particular category in your configuration files. This can be useful when you need to understand the structure of a specific section for easier navigation and management.

```sql
select
  path,
  section,
  comment
from
  ini_section
where
  section like 'settings.%';
```