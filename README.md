# nippo-cli

<!-- markdownlint-disable MD013 -->

![Coverage](https://raw.githubusercontent.com/c18t/nippo-cli/main/badges/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/c18t/nippo-cli/main/badges/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/c18t/nippo-cli/main/badges/time.svg)

<!-- markdownlint-enable MD013 -->

The tool to power my nippo.

## Install `nippo` command

```shell
go install github.com/c18t/nippo-cli/nippo@latest
```

## Usage

### Setup

1. Initialize the configuration:

   ```shell
   nippo init
   ```

   This will prompt you for:
   - Google Drive folder URL/ID
   - Site URL
   - Project URL
   - Project branch name
   - Template and asset paths

2. Download `credentials.json` from Google Cloud Console:
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create or select a project
   - Enable the Google Drive API
   - Create OAuth 2.0 credentials (Desktop application)
   - Download the credentials and save as `credentials.json` in the data directory

3. Authenticate with Google Drive:

   ```shell
   nippo auth
   ```

4. Check your setup:

   ```shell
   nippo doctor
   ```

### Build

```shell
nippo build
```

### Publish

```shell
nippo deploy
```

## Configuration

### Configuration File

The configuration file is located at:

| Platform    | Path                         |
| ----------- | ---------------------------- |
| Linux/macOS | `~/.config/nippo/nippo.toml` |
| Windows     | `%APPDATA%\nippo\nippo.toml` |

You can also set `XDG_CONFIG_HOME` environment variable to customize the location.

### Path Configuration

nippo uses XDG Base Directory specification for file locations.
You can customize paths in `nippo.toml`:

```toml
[project]
drive_folder_id = "your-drive-folder-id"
site_url = "https://nippo.example.com"
url = "https://github.com/c18t/nippo"
branch = "main"
template_path = "/templates"
asset_path = "/dist"

[path]
# Uncomment and modify to customize file locations.
# data_dir = "~/.local/share/nippo"
# cache_dir = "~/.cache/nippo"
```

### Default Paths

#### Data Directory

Files: `credentials.json`, `token.json`, `templates/`, `assets/`

| Platform    | Default Path                                     |
| ----------- | ------------------------------------------------ |
| Linux/macOS | `$XDG_DATA_HOME/nippo` or `~/.local/share/nippo` |
| Windows     | `%LOCALAPPDATA%\nippo`                           |

#### Cache Directory

Files: `md/`, `output/`, `nippo-template.zip`

| Platform    | Default Path                                |
| ----------- | ------------------------------------------- |
| Linux/macOS | `$XDG_CACHE_HOME/nippo` or `~/.cache/nippo` |
| Windows     | `%LOCALAPPDATA%\nippo\cache`                |

### Example: Custom Paths

```toml
[project]
drive_folder_id = "1HNSRS2tJI2t7DKP_8XQJ2NTleSH-rs4y"
site_url = "https://nippo.c18t.me"
url = "https://github.com/c18t/nippo"
branch = "main"
template_path = "/templates"
asset_path = "/dist"

[path]
# Store data in a custom location
data_dir = "/opt/nippo/data"
cache_dir = "/var/cache/nippo"
```

You can also use:

- Environment variables: `$HOME/nippo-data` or `%USERPROFILE%\nippo-data`
- Tilde expansion: `~/nippo-data`
- Relative paths: `./data` (relative to config directory)

## Development

### Getting Started

#### Prerequisites

- [Codespaces](https://github.co.jp/features/codespaces)

Or some IDE with [Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers)
support (e.g., Visual Studio Code).

#### Create Container

##### Using GitHub Codespaces

1. Open this repository in GitHub
2. Click "Code" → "Codespaces" → "Create codespace on main"
3. Wait for the Codespace to be created and initialized
4. Once set up, proceed to the "Setup Container Workspace" section below

##### Using Dev Containers (Local)

1. Clone the repository:

   ```shell
   ghq get c18t/nippo-cli
   cd $(ghq root)/github.com/c18t/nippo-cli
   ```

2. Add GH_TOKEN to .env (if necessary):

   ```shell
   cp .env.sample .env
   gh auth token | xargs -I {} echo "GH_TOKEN="{} >> .env
   ```

3. Open the project in Dev Containers:
   1. `code .`
   2. `Ctrl` + `Shift` + `P`
   3. `>Dev Containers: Reopen in Container`

#### Setup Container Workspace

1. Run setup tasks:

   ```shell
   post-create.sh
   ```

2. Build and run the application:

   ```shell
   mise run build
   ./bin/nippo
   ```

3. [extra] Install extensions recommended for the workspace:
   1. `Ctrl` + `Shift` + `P`
   2. `>Extensions: Show Recommended Extensions`
   3. Click `install` button

### Available Task Runner Commands

`mise run <task name>`

```console
$ mise tasks
Name                                 Description
build                                Build the CLI application
dev-up:ccmanager-worktree-settings   Set up ccmanager worktree auto-directory settings
dev-up:claude-code-skip-permissions  Set up claude code to skip permissions
dev-up:claude-code-stop-autoupdates  Set up claude code to disable auto-updates
devcontainer-up                      Start devcontainer and run ccmanager with Claude Code
release                              Build release binaries
setup                                Set up (Runs all `setup:*` tasks)
setup:claude-mcp                     Set up Claude Code MCP servers
setup:go-mod                         Install go modules with go.mod
setup:ignore-workspace-file-changes  Ignore local changes to workspace file
setup:mise                           Install dev dependencies with mise
setup:pnpm                           Set up pnpm packages
setup:pre-commit                     Set up pre-commit hooks
```

See Also: [c18t/boilerplate-go-cli](https://github.com/c18t/boilerplate-go-cli)

## License

[MIT](./LICENSE)

## Author

ɯ̹t͡ɕʲi

- [github / c18t](https://github.com/c18t)
