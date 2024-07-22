# Changelog

## [v0.14.1](https://github.com/c18t/nippo-cli/compare/v0.14.0...v0.14.1) - 2024-07-22

- add github action and pr template by @c18t in https://github.com/c18t/nippo-cli/pull/3

## [v0.14.0](https://github.com/c18t/nippo-cli/compare/v0.13.2...v0.14.0) - 2024-01-03

- update init command to add project settings and interactive UI

## [v0.13.2](https://github.com/c18t/nippo-cli/compare/v0.13.1...v0.13.2) - 2024-01-01

- add meta description to templates
  - Include meta description in templates for improved SEO.

## [v0.13.1](https://github.com/c18t/nippo-cli/compare/v0.13.0...v0.13.1) - 2023-12-31

- update build command
  - Increase the number of items in feed.xml from 10 to 20.

## [v0.13.0](https://github.com/c18t/nippo-cli/compare/v0.12.1...v0.13.0) - 2023-12-31

- update build command
  - Add sitemap generation feature.

## [v0.12.1](https://github.com/c18t/nippo-cli/compare/v0.12.0...v0.12.1) - 2023-12-31

- fix vercel deploy command
  - Compress into single file before uploading output.

## [v0.12.0](https://github.com/c18t/nippo-cli/compare/v0.11.0...v0.12.0) - 2023-12-31

- update build command
  - Add article feed generation feature
  - Add canonical url link to templates

## [v0.11.0](https://github.com/c18t/nippo-cli/compare/v0.10.1...v0.11.0) - 2023-12-31

- update build command
  - Add a feature to recursively download the nippo folder, including subfolders.
  - Add a feature to support monthly archive page building for all months, not just the current month.

## [v0.10.1](https://github.com/c18t/nippo-cli/compare/v0.10.0...v0.10.1) - 2023-12-31

- change nippo data file name format
  - Change file name format from yyyymmdd to yyyy-mm-dd.
  - Allow any characters after the date in file name.
  - [WIP] Add a feature to recursively download the nippo folder.

## [v0.10.0](https://github.com/c18t/nippo-cli/compare/v0.9.0...v0.10.0) - 2023-12-30

- update build command
  - Add support for incremental download of nippo data.
- refactor codes.

## [v0.9.0](https://github.com/c18t/nippo-cli/compare/v0.8.0...v0.9.0) - 2023-12-25

- update deploy command to support asset deployment

## [v0.8.0](https://github.com/c18t/nippo-cli/compare/v0.7.0...v0.8.0) - 2023-12-25

- update build command to support OGP

## [v0.7.0](https://github.com/c18t/nippo-cli/compare/v0.6.0...v0.7.0) - 2023-12-25

- update clean command to add file deletion features

## [v0.6.0](https://github.com/c18t/nippo-cli/compare/v0.5.2...v0.6.0) - 2023-12-24

- update build command to add archive page build feature

## [v0.5.2](https://github.com/c18t/nippo-cli/compare/v0.5.1...v0.5.2) - 2023-12-24

- add path and title string output methods to NippoDate

## [v0.5.1](https://github.com/c18t/nippo-cli/compare/v0.5.0...v0.5.1) - 2023-12-24

- fix deploy command: os/exec package usage error

## [v0.5.0](https://github.com/c18t/nippo-cli/compare/v0.4.2...v0.5.0) - 2023-12-24

- update update and deploy commands

## [v0.4.2](https://github.com/c18t/nippo-cli/compare/v0.4.1...v0.4.2) - 2023-12-24

- add build artifact cleanup to build command

## [v0.4.1](https://github.com/c18t/nippo-cli/compare/v0.4.0...v0.4.1) - 2023-12-24

- fix an issue where the `clean` controller is called for `init`, `update`, and `deploy` commands.

## [v0.4.0](https://github.com/c18t/nippo-cli/compare/v0.3.2...v0.4.0) - 2023-12-23

- update build command to add nippo page generating

## [v0.3.2](https://github.com/c18t/nippo-cli/compare/v0.3.1...v0.3.2) - 2023-12-23

- refactor dependency management based on Clean Architecture

## [v0.3.1](https://github.com/c18t/nippo-cli/compare/v0.3.0...v0.3.1) - 2023-12-17

- refactor code about config and template logics

## [v0.3.0](https://github.com/c18t/nippo-cli/compare/v0.3.0...v0.3.1) - 2023-12-17

- update init and build command to add top page generating

## [v0.2.0](https://github.com/c18t/nippo-cli/compare/v0.1.0...v0.2.0) - 2023-12-17

- update init and build command to add oauth/download features

## [v0.1.0](https://github.com/c18t/nippo-cli/compare/v0.0.1...v0.1.0) - 2023-12-14

- update init command to add download feature for c18t/nippo project
- update docker-compose settings for devcontainer
- add devcontainer and vscode debug config

## [v0.0.1](https://github.com/c18t/nippo-cli/compare/0a388100d49db6775647808ab6cba61cd2cd029e...v0.0.1) - 2023-12-11

- add subcommands
