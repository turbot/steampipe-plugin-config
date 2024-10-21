## v1.0.0 [2024-10-22]

There are intentionally no significant changes in this plugin version, but it has been released to coincide with the [Steampipe's v1.0.0](https://steampipe.io/changelog/steampipe-cli-v1-0-0) release. This plugin follows [semantic versioning's specification](https://semver.org/#semantic-versioning-specification-semver) and preserves backward compatibility in each major version.

_Dependencies_

- Recompiled plugin with Go version `1.22`. ([#57](https://github.com/turbot/steampipe-plugin-config/pull/57))
- Recompiled plugin with [steampipe-plugin-sdk v5.10.4](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5104-2024-08-29) that fixes logging in the plugin export tool. ([#57](https://github.com/turbot/steampipe-plugin-config/pull/57))

## v0.6.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#52](https://github.com/turbot/steampipe-plugin-config/pull/52))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#52](https://github.com/turbot/steampipe-plugin-config/pull/52))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-config/blob/main/docs/LICENSE). ([#52](https://github.com/turbot/steampipe-plugin-config/pull/52))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#51](https://github.com/turbot/steampipe-plugin-config/pull/51))

## v0.5.1 [2023-10-05]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#44](https://github.com/turbot/steampipe-plugin-config/pull/44))

## v0.5.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#42](https://github.com/turbot/steampipe-plugin-config/pull/42))
- Recompiled plugin with Go version `1.21`. ([#42](https://github.com/turbot/steampipe-plugin-config/pull/42))

## v0.4.0 [2023-04-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#36](https://github.com/turbot/steampipe-plugin-config/pull/36))

## v0.3.1 [2023-02-28]

_Bug fixes_

- Fixed the incorrect dynamic plugin schema definition that would cause 100% CPU during initialization. ([#34](https://github.com/turbot/steampipe-plugin-config/pull/34))

## v0.3.0 [2022-11-17]

_What's new?_

- Added support for retrieving INI, JSON and YML config files from remote Git repositories and S3 buckets. For more information, please see [Supported Path Formats](https://hub.steampipe.io/plugins/turbot/config#supported-path-formats). ([#27](https://github.com/turbot/steampipe-plugin-config/pull/27))
- Added file watching support for files included in the `paths` config argument. ([#27](https://github.com/turbot/steampipe-plugin-config/pull/27))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.0.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v500-2022-11-16) which includes support for fetching remote files with go-getter and file watching. ([#27](https://github.com/turbot/steampipe-plugin-config/pull/27))

## v0.2.0 [2022-09-26]

_Enhancements_

- Added support for loading JSON files containing arrays to `json_file` table. ([#25](https://github.com/turbot/steampipe-plugin-config/pull/25)) (Thanks to [@keilin-anz](https://github.com/keilin-anz) for the contribution!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#26](https://github.com/turbot/steampipe-plugin-config/pull/26))
- Recompiled plugin with Go version `1.19`. ([#26](https://github.com/turbot/steampipe-plugin-config/pull/26))

## v0.1.0 [2022-04-27]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#19](https://github.com/turbot/steampipe-plugin-config/pull/19))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#20](https://github.com/turbot/steampipe-plugin-config/pull/20))

## v0.0.3 [2022-03-17]

_What's new?_

- New tables added
  - [json_key_value](https://hub.steampipe.io/plugins/turbot/config/tables/json_key_value) ([#12](https://github.com/turbot/steampipe-plugin-config/pull/12))
  - [yml_key_value](https://hub.steampipe.io/plugins/turbot/config/tables/yml_key_value) ([#12](https://github.com/turbot/steampipe-plugin-config/pull/12))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.0.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v301-2022-03-10) ([#12](https://github.com/turbot/steampipe-plugin-config/pull/12))

## v0.0.2 [2022-02-14]

_Enhancements_

- Update examples in README, index doc, and table docs for better rendering in the Hub ([#16](https://github.com/turbot/steampipe-plugin-config/pull/16))

## v0.0.1 [2022-02-11]

_What's new?_

- New tables added
  - [ini_key_value](https://hub.steampipe.io/plugins/turbot/config/tables/ini_key_value)
  - [ini_section](https://hub.steampipe.io/plugins/turbot/config/tables/ini_section)
  - [json_file](https://hub.steampipe.io/plugins/turbot/config/tables/json_file)
  - [yml_file](https://hub.steampipe.io/plugins/turbot/config/tables/yml_file)
