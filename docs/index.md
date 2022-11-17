---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/config.svg"
brand_color: "#222017"
display_name: "Config"
short_name: "config"
description: "Steampipe plugin to query data from various types of files like INI, JSON, YML and more."
og_description: "Query data from various types of files with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/config-social-graphic.png"
---

# Config + Steampipe

Config plugin is used to parse various types of configuration files, e.g., `INI`, `JSON`, `YML`, in order to represent the content as SQL tables.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query data using SQL.

The following file types are currently supported:

- INI
- JSON
- YML

Query all data in your INI files:

```sql
select
  path,
  section,
  key,
  value
from
  ini_key_value;
```

```sh
+----------------------------+----------+---------------+-------------------------------------------+
| path                       | section  | key           | value                                     |
+----------------------------+----------+---------------+-------------------------------------------+
| /Users/myuser/defaults.ini | Settings | DetailedLog   | 1                                         |
| /Users/myuser/defaults.ini | Status   | RunStatus     | 1                                         |
| /Users/myuser/defaults.ini | Status   | StatusRefresh | 10                                        |
| /Users/myuser/defaults.ini | Status   | StatusPort    | 6090                                      |
| /Users/myuser/logs.ini     | Server   | Archive       | 1                                         |
| /Users/myuser/logs.ini     | Server   | ServerName    | Unknown                                   |
| /Users/myuser/logs.ini     | Settings | LogFile       | /opt/ecs/mvuser/MV_IPTel/log/MV_IPTel.log |
| /Users/myuser/logs.ini     | Settings | Version       | 0.9 Build 4 Created July 11 2004 14:00    |
+----------------------------+----------+---------------+-------------------------------------------+
```

Query all data in your JSON files:

```sql
select
  path,
  jsonb_pretty(content) as file_content
from
  json_file;
```

```sh
+----------------------------+------------------------------------------------------------+
| path                       | file_content                                               |
+----------------------------+------------------------------------------------------------+
| /Users/myuser/invoice.json | {                                                          |
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
|                            |     "bill-to": null,                                       |
|                            |     "receipt": "Oz-Ware Purchase Invoice",                 |
|                            |     "ship-to": null,                                       |
|                            |     "customer": {                                          |
|                            |         "first_name": "Dorothy",                           |
|                            |         "family_name": "Gale"                              |
|                            |     },                                                     |
|                            | }                                                          |
| /Users/myuser/test.json    | {                                                          |
|                            |     "foo": "bar",                                          |
|                            |     "includes": [                                          |
|                            |         "common.json"                                      |
|                            |     ]                                                      |
|                            | }                                                          |
+----------------------------+------------------------------------------------------------+
```

Query all data in your YML files:

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

## Documentation

- **[Table definitions & examples →](/plugins/turbot/config/tables)**

## Get started

### Install

Download and install the latest Config plugin:

```bash
steampipe plugin install config
```

### Credentials

No credentials are required.

### Configuration

Installing the latest config plugin will create a config file (`~/.steampipe/config/config.spc`) with a single connection named `config`:

```hcl
connection "config" {
  plugin = "config"

  # Each paths argument is a list of locations to search for a particular file type
  # Each paths can be configured with a local directory, a remote Git repository URL, or an S3 bucket URL
  # Wildcard based searches are supported, including recursive searches
  # Local paths are resolved relative to the current working directory (CWD)

  # For example, for the json_paths argument:
  #  - "*.json" matches all JSON files in the CWD
  #  - "**/*.json" matches all JSON files in a directory, and all the sub-directories in it
  #  - "../*.json" matches all JSON files in in the CWD's parent directory
  #  - "steampipe*.json" matches all JSON files starting with "steampipe" in the CWD
  #  - "/path/to/dir/*.json" matches all JSON files in a specific directory
  #  - "/path/to/dir/main.json" matches a specific file

  # If paths includes "*", all files (including non-required files) in
  # the CWD will be matched, which may cause errors if incompatible file types exist

  # All paths arguments defaults to CWD
  ini_paths  = [ "*.ini" ]
  json_paths = [ "*.json" ]
  yml_paths  = [ "*.yml", "*.yaml" ]
}
```

### Supported Path Formats

The `ini_paths`, `json_paths` and `yml_paths` config arguments are flexible and can search for INI, JSON and YML files from several different sources respectively, e.g., local directory paths, Git, S3.

The following sources are supported:

- [Local files](#configuring-local-file-paths)
- [Remote Git repositories](#configuring-remote-git-repository-urls)
- [S3](#configuring-s3-urls)

Paths may [include wildcards](https://pkg.go.dev/path/filepath#Match) and support `**` for recursive matching. For example:

```hcl
connection "config" {
  plugin = "config"

  ini_paths = [
    "*.ini",
    "~/*.ini",
    "github.com/madmurphy/libconfini//examples/ini_files//l*.ini",
    "github.com/madmurphy/libconfini//examples/ini_files//**/*.ini",
    "bitbucket.org/feeeper/ini-examples//**/*.ini",
    "s3::https://bucket.s3.amazonaws.com/config_examples//**/*.ini"
  ]

  json_paths = [
    "*.json",
    "~/*.json",
    "github.com/LearnWebCode/json-example//*.json",
    "github.com/LearnWebCode/json-example//**/*.json",
    "bitbucket.org/atlassian/json-schema-diff//**/*.json",
    "s3::https://bucket.s3.amazonaws.com/config_examples//**/*.json"
  ]

  yml_paths  = [
    "*.yml",
    "~/*.yaml",
    "github.com/awslabs/aws-cloudformation-templates//aws/services/ElasticLoadBalancing//*.yaml",
    "github.com/awslabs/aws-cloudformation-templates//**/*.yaml",
    "bitbucket.org/lokerd/yaml-formation//**/*.yml",
    "gitlab.com/versioncontrol1/cloudformationproject//substacks//*.yaml",
    "s3::https://bucket.s3.amazonaws.com/config_examples//**/*.yaml"
  ]
}
```

**Note**: If any path matches on `*` without the expected file extension (e.g. `.ini`, `.json`, `.yaml`, `.yml`), all files (including non-required files) in the directory will be matched, which may cause errors if incompatible file types exist.

#### Configuring Local File Paths

You can define a list of local directory paths to search for INI, JSON and YML files. Paths are resolved relative to the current working directory.

For example, for the `json_paths` argument:

- `*.json` matches all JSON files in the CWD.
- `**/*.json` matches all JSON files in the CWD and all sub-directories.
- `../*.json` matches all JSON files in the CWD's parent directory.
- `steampipe*.json` matches all JSON files starting with `steampipe` in the CWD.
- `/path/to/dir/*.json` matches all JSON files in a specific directory. For example:
  - `~/*.json` matches all JSON files in the home directory.
  - `~/**/*.json` matches all JSON files recursively in the home directory.
- `/path/to/dir/main.json` matches a specific file.

```hcl
connection "config" {
  plugin = "config"

  ini_paths  = [ "*.ini", "~/*.ini", "/path/to/dir/main.ini" ]
  json_paths = [ "*.json", "~/*.json", "/path/to/dir/main.json" ]
  yml_paths  = [ "*.yml", "~/*.yaml", "/path/to/dir/main.yml" ]
}
```

**NOTE:** If paths includes `*`, all files (including non-required files) in the CWD will be matched, which may cause errors if incompatible file types exist.

#### Configuring Remote Git Repository URLs

You can also configure `ini_paths`, `json_paths` and `yml_paths` with any Git remote repository URLs, e.g., GitHub, BitBucket, GitLab. The plugin will then attempt to retrieve any INI, JSON and YML files from the remote repositories.

For example:

- `github.com/LearnWebCode/json-example//*.json` matches all top-level JSON files in the specified github repository.
- `github.com/LearnWebCode/json-example//**/*.json` matches all JSON files in the specified github repository and all sub-directories.
- `github.com/turbot/polygoat?ref=fix_7677//**/*.json` matches all JSON files in the specific tag of github repository.
- `github.com/brandiqa/json-examples//api/config//*.json` matches all JSON files in the specified folder path.

You can specify a subdirectory after a double-slash (`//`) if you want to download only a specific subdirectory from a downloaded directory.

```hcl
connection "config" {
  plugin = "config"

  paths = [ "github.com/brandiqa/json-examples//api/config//*.json" ]
}
```

Similarly, you can define a list of GitLab and BitBucket URLs to search for INI, JSON and YML files:

```hcl
connection "config" {
  plugin = "config"

  ini_paths  = [
    "github.com/madmurphy/libconfini//examples/ini_files//l*.ini",
    "github.com/madmurphy/libconfini//**/*.ini",
    "bitbucket.org/feeeper/ini-examples//**/*.ini",
    "bitbucket.org/feeeper/ini-examples//*.ini"
  ]

  json_paths = [
    "github.com/LearnWebCode/json-example//*.json",
    "github.com/brandiqa/json-examples//api/config//*.json",
    "bitbucket.org/atlassian/json-schema-diff//**/*.json",
    "bitbucket.org/atlassian/json-schema-diff//*.json"
  ]

  yml_paths  = [
    "github.com/awslabs/aws-cloudformation-templates//aws/services/ElasticLoadBalancing//*.yaml",
    "github.com/awslabs/aws-cloudformation-templates//aws/services//**/*.yaml",
    "bitbucket.org/lokerd/yaml-formation//**/*.yml",
    "bitbucket.org/lokerd/yaml-formation//examples/convert/Parameters//*.yml",
    "gitlab.com/versioncontrol1/cloudformationproject//**/*.yaml",
    "gitlab.com/versioncontrol1/cloudformationproject//substacks//*.yaml"
  ]
}
```

#### Configuring S3 URLs

You can also query all INI, JSON and YML files stored inside an S3 bucket (public or private) using the bucket URL.

##### Accessing a Private Bucket

In order to access your files in a private S3 bucket, you will need to configure your credentials. You can use your configured AWS profile from local `~/.aws/config`, or pass the credentials using the standard AWS environment variables, e.g., `AWS_PROFILE`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`.

We recommend using AWS profiles for authentication.

**Note:** Make sure that `region` is configured in the config. If not set in the config, `region` will be fetched from the standard environment variable `AWS_REGION`.

You can also authenticate your request by setting the AWS profile and region in the respective path argument. For example:

```hcl
connection "config" {
  plugin = "config"

  paths = [
    "s3::https://bucket-2.s3.us-east-1.amazonaws.com//*.json?aws_profile=<AWS_PROFILE>",
    "s3::https://bucket-2.s3.us-east-1.amazonaws.com/test_folder//*.json?aws_profile=<AWS_PROFILE>"
  ]
}
```

**Note:**

In order to access the bucket, the IAM user or role will require the following IAM permissions:

- `s3:ListBucket`
- `s3:GetObject`
- `s3:GetObjectVersion`

If the bucket is in another AWS account, the bucket policy will need to grant access to your user or role. For example:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ReadBucketObject",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::123456789012:user/YOUR_USER"
      },
      "Action": ["s3:ListBucket", "s3:GetObject", "s3:GetObjectVersion"],
      "Resource": ["arn:aws:s3:::test-bucket1", "arn:aws:s3:::test-bucket1/*"]
    }
  ]
}
```

##### Accessing a Public Bucket

Public access granted to buckets and objects through ACLs and bucket policies allows any user access to data in the bucket. We do not recommend making S3 buckets public, but if there are specific objects you'd like to make public, please see [How can I grant public read access to some objects in my Amazon S3 bucket?](https://aws.amazon.com/premiumsupport/knowledge-center/read-access-objects-s3-bucket/).

You can query any public S3 bucket directly using the URL without passing credentials. For example:

```hcl
connection "config" {
  plugin = "config"

  paths = [
    "s3::https://bucket-1.s3.us-east-1.amazonaws.com/test_folder//*.json",
    "s3::https://bucket-2.s3.us-east-1.amazonaws.com/test_folder//**/*.json"
  ]
}
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-config
- Community: [Slack Channel](https://steampipe.io/community/join)
