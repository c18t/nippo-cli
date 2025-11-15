#!/bin/bash
set -euo pipefail  # Exit on error, undefined vars, and pipeline failures

USER=$(whoami)

# mise settings
# shellcheck disable=SC2016
echo 'eval "$(mise activate bash)"' >> "/home/${USER}/.bashrc"
mise trust
mise run setup:mise

# pnpm settings
pnpm setup
export PNPM_HOME="/home/${USER}/.local/share/pnpm"
case ":$PATH:" in
  *":$PNPM_HOME:"*) ;;
  *) export PATH="$PNPM_HOME:$PATH" ;;
esac
# 明示的に指定しないとワークスペースに作られてしまう
pnpm config set store-dir "${PNPM_HOME}/store"
# node-gyp がないと node-pty のビルドスクリプトで落ちる
pnpm i -g node-gyp
# mise run setup:pnpm で初回インストールすると非常に古い pnpm で動いてしまうため個別に実行
pnpm i

mise run setup
