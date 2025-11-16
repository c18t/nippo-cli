# nippo-cli

The tool to power my nippo.

## Install `nippo` command

```shell
go install github.com/c18t/nippo-cli/nippo@latest
```

## Usage

### Setup

```shell
nippo init
```

### Build

```shell
nippo build
```

### Publish

```shell
nippo deploy
```

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
dev-up:ccmanager-skip-permissions    Set up ccmanager to skip permissions
dev-up:ccmanager-worktree-settings   Set up ccmanager worktree auto-directory settings
dev-up:claude-code-stop-autoupdates  Set up Claude Code to disable auto-updates
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
