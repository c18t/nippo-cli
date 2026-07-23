# Changelog

## [v0.16.6](https://github.com/c18t/nippo-cli/compare/v0.16.5...v0.16.6) - 2026-07-23

- Release for dependabot updates (#119, #120, #121) by @c18t in https://github.com/c18t/nippo-cli/pull/124
- chore(deps): bump google.golang.org/api from 0.271.0 to 0.288.0 in the gomod-minor-patch group across 1 directory by @dependabot[bot] in https://github.com/c18t/nippo-cli/pull/119
- chore(deps): bump actions/checkout from 6 to 7 in the github-actions group by @dependabot[bot] in https://github.com/c18t/nippo-cli/pull/120
- chore(deps-dev): bump the npm-dev-dependencies group with 2 updates by @dependabot[bot] in https://github.com/c18t/nippo-cli/pull/121

## [v0.16.5](https://github.com/c18t/nippo-cli/compare/v0.16.4...v0.16.5) - 2026-07-12

- docs(github): add pull request template for release integration PRs by @c18t in https://github.com/c18t/nippo-cli/pull/113
- docs(github): improve pull request templates by @c18t in https://github.com/c18t/nippo-cli/pull/116
- Release: DX improvements (#106-#111) by @c18t in https://github.com/c18t/nippo-cli/pull/115
- chore(pnpm): limit network concurrency to avoid socket timeouts by @c18t in https://github.com/c18t/nippo-cli/pull/106
- ci(octocov): publish badges to dedicated badges branch by @c18t in https://github.com/c18t/nippo-cli/pull/107
- chore: fix always-failing pre-commit hooks (govulncheck, bearer) by @c18t in https://github.com/c18t/nippo-cli/pull/108
- ci(dependabot): group updates to avoid lockfile conflict chains by @c18t in https://github.com/c18t/nippo-cli/pull/109
- chore(lint): use pnpm-managed markdownlint-cli2 in pre-commit by @c18t in https://github.com/c18t/nippo-cli/pull/110
- docs(openspec): stop pinning go version in project.md by @c18t in https://github.com/c18t/nippo-cli/pull/111

## [v0.16.4](https://github.com/c18t/nippo-cli/compare/v0.16.3...v0.16.4) - 2026-07-12

- chore(deps): clean up npm dev dependencies and fix octocov on dependabot runs by @c18t in https://github.com/c18t/nippo-cli/pull/103

## [v0.16.3](https://github.com/c18t/nippo-cli/compare/v0.16.2...v0.16.3) - 2026-07-11

- Bump @fission-ai/openspec from 0.16.0 to 1.6.0 by @dependabot[bot] in https://github.com/c18t/nippo-cli/pull/82

## [v0.16.2](https://github.com/c18t/nippo-cli/compare/v0.16.1...v0.16.2) - 2026-07-11

- Bump google.golang.org/api from 0.257.0 to 0.271.0 by @dependabot[bot] in https://github.com/c18t/nippo-cli/pull/100
- Bump goreleaser/goreleaser-action from 6 to 7 by @dependabot[bot] in https://github.com/c18t/nippo-cli/pull/97
- Bump markdownlint-cli2 from 0.19.1 to 0.20.0 by @dependabot[bot] in https://github.com/c18t/nippo-cli/pull/85
- Bump github.com/charmbracelet/bubbles from 0.21.0 to 1.0.0 by @dependabot[bot] in https://github.com/c18t/nippo-cli/pull/95

## [v0.16.1](https://github.com/c18t/nippo-cli/compare/v0.16.0...v0.16.1) - 2025-12-14

- fix(presenter): resolve nil pointer panic on version command by @c18t in https://github.com/c18t/nippo-cli/pull/78
- fix issue where user's ~/.config/nippo/nippo.toml is overwritten during test execution by @c18t in https://github.com/c18t/nippo-cli/pull/78

## [v0.16.0](https://github.com/c18t/nippo-cli/compare/v0.15.4...v0.16.0) - 2025-12-14

- Release v0.16.0: TUI components, front-matter support, auth/doctor commands, and security fixes by @c18t in https://github.com/c18t/nippo-cli/pull/73

BRAKING CHANGES:

- OAuth scope has been changed to DriveScope (write permission), so existing users need to re-authenticate with `nippo auth` command.

## [v0.15.4](https://github.com/c18t/nippo-cli/compare/v0.15.3...v0.15.4) - 2025-12-12

- chore(deps): update dependencies by @c18t in https://github.com/c18t/nippo-cli/pull/69
- Bump github.com/spf13/cobra from 1.10.1 to 1.10.2 by @dependabot[bot] in https://github.com/c18t/nippo-cli/pull/68

## [v0.15.3](https://github.com/c18t/nippo-cli/compare/v0.15.2...v0.15.3) - 2025-11-17

- chore: update devcontainer and dependencies by @c18t in https://github.com/c18t/nippo-cli/pull/59

## [v0.15.2](https://github.com/c18t/nippo-cli/compare/v0.15.1...v0.15.2) - 2025-11-17

- fix(release): resolve GoReleaser configuration issues by @c18t in https://github.com/c18t/nippo-cli/pull/57

## [v0.15.1](https://github.com/c18t/nippo-cli/compare/v0.15.0...v0.15.1) - 2025-11-16

- Fix goreleaser config and markdownlint errors by @c18t in https://github.com/c18t/nippo-cli/pull/55

## [v0.15.0](https://github.com/c18t/nippo-cli/compare/v0.14.4...v0.15.0) - 2025-11-16

- boilerplate-go-cliのアップデートを取り込む by @c18t in https://github.com/c18t/nippo-cli/pull/42
- chore(devenv): update development environment from boilerplate-go-cli by @c18t in https://github.com/c18t/nippo-cli/pull/52
- Refactor DI implementation with package pattern and lifecycle management by @c18t in https://github.com/c18t/nippo-cli/pull/53

## [v0.14.4](https://github.com/c18t/nippo-cli/compare/v0.14.3...v0.14.4) - 2024-07-22

- replace DI package `dig` with `do/v2` by @c18t in https://github.com/c18t/nippo-cli/pull/19

## [v0.14.3](https://github.com/c18t/nippo-cli/compare/v0.14.2...v0.14.3) - 2024-07-22

- add dependabot settings by @c18t in https://github.com/c18t/nippo-cli/pull/10
- Bump github.com/spf13/cobra from 1.8.0 to 1.8.1 by @dependabot in https://github.com/c18t/nippo-cli/pull/13
- Bump github.com/gorilla/feeds from 1.1.2 to 1.2.0 by @dependabot in https://github.com/c18t/nippo-cli/pull/12
- Bump google.golang.org/api from 0.154.0 to 0.189.0 by @dependabot in https://github.com/c18t/nippo-cli/pull/14
- Bump github.com/spf13/viper from 1.18.1 to 1.19.0 by @dependabot in https://github.com/c18t/nippo-cli/pull/16

## [v0.14.3](https://github.com/c18t/nippo-cli/compare/v0.14.2...v0.14.3) - 2024-07-22

- add dependabot settings by @c18t in https://github.com/c18t/nippo-cli/pull/10

## [v0.14.2](https://github.com/c18t/nippo-cli/compare/v0.14.1...v0.14.2) - 2024-07-22

- align devcontainer.json and compose.yaml configurations by @c18t in https://github.com/c18t/nippo-cli/pull/6

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
