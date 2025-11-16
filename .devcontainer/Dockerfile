# syntax=docker/dockerfile:1
# hadolint global ignore=DL3008
FROM debian:trixie-slim

# ------------------------
# devcontainer information
# ------------------------
# set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

# workspace settings
ARG WORKSPACE_ROOT="/workspaces"
ARG WORKING_DIR="${WORKSPACE_ROOT}/app"

# container user settings
ARG CONTAINER_USER=user
ARG UID=1000
ARG GID=1000
ARG TZ="Asia/Tokyo"

# Delve debug server port
ARG DEBUG_PORT_INTERNAL=2345
EXPOSE ${DEBUG_PORT_INTERNAL}

# Mountpoint for cache
ENV COMMANDHISTORY="/commandhistory"
ENV MISEPATH="/mise"
ENV PNPM_HOME="/home/${CONTAINER_USER}/.local/share/pnpm"
ENV NODE_MODULES="${WORKING_DIR}/node_modules"
ENV CLAUDE_CODE_CONFIG="/home/${CONTAINER_USER}/.claude"

# running process
CMD [ "sleep", "infinity" ]

# -----------
# build stage
# -----------
SHELL ["/bin/bash", "-eo", "pipefail", "-c"]

# shell settings
ENV SHELL=bash
ENV LANG="C.UTF-8"
ENV TMP=/tmp
ENV TEMP=${TMP}

# mise settings
ENV MISE_DATA_DIR="${MISEPATH}"
ENV MISE_CONFIG_DIR="${MISEPATH}"
ENV MISE_CACHE_DIR="${MISEPATH}/cache"
ENV MISE_INSTALL_PATH="/usr/local/bin/mise"
ENV PATH="${MISEPATH}/shims:$PATH"

ARG XDG_LOCAL="/home/${CONTAINER_USER}/.local"
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    bash <<EOF
set -eo pipefail

# Add Debian sid repository to install the latest git version
# This ensures access to the latest git features and security patches
# Pin priority (990) ensures git packages are always pulled from sid
echo 'deb http://deb.debian.org/debian sid main' > /etc/apt/sources.list.d/sid.list
printf 'Package: *\nPin: release a=sid\nPin-Priority: 100\n\n' > /etc/apt/preferences.d/git-from-sid
printf 'Package: git*\nPin: release a=sid\nPin-Priority: 990\n' >> /etc/apt/preferences.d/git-from-sid

# Install development tools and utilities
# cf. https://mise.jdx.dev/mise-cookbook/docker.html
#     https://code.claude.com/docs/en/devcontainer
apt-get update
apt-get -y --no-install-recommends install \
  sudo curl git ca-certificates build-essential `# Core development tools and mise dependencies`\
  less procps man-db unzip gnupg2 `# Debugging and package management utilities` \
  iptables ipset iproute2 dnsutils aggregate jq `# Networking tools for init-firewall`

# Install mise (polyglot runtime manager)
curl https://mise.run | sh

# Create non-root container user
groupadd -g \${GID} -o \${CONTAINER_USER}
useradd -m -s /bin/bash -u \${UID} -g \${GID} -o \${CONTAINER_USER}
echo "\${CONTAINER_USER} ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/\${CONTAINER_USER}
chmod 0440 /etc/sudoers.d/\${CONTAINER_USER}

# Create workspace, mise data, and configuration directories
mkdir -p \
  \${MISEPATH} \
  \${XDG_LOCAL} \
  \${PNPM_HOME} \
  \${WORKSPACE_ROOT} \
  \${WORKING_DIR} \
  \${WORKING_DIR}/node_modules \
  \${CLAUDE_CODE_CONFIG}

# Grant ownership of all configuration directories to the container user
chown -R \${CONTAINER_USER} \
  \${MISEPATH} \
  \${XDG_LOCAL} \
  \${PNPM_HOME} \
  \${WORKSPACE_ROOT} \
  \${WORKING_DIR} \
  \${WORKING_DIR}/node_modules \
  \${CLAUDE_CODE_CONFIG}

# Configure persistent bash history across container restarts
SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=\${COMMANDHISTORY}/.bash_history"
mkdir -p \${COMMANDHISTORY}
touch \${COMMANDHISTORY}/.bash_history
chown -R \${CONTAINER_USER} \${COMMANDHISTORY}
echo "\${SNIPPET}" >> "/home/\${CONTAINER_USER}/.bashrc"
EOF

# Copy project-specific scripts
COPY --chmod=755 start.sh post-create.sh init-firewall.sh /usr/local/bin/

# skip firewall initialization
# RUN echo "%${CONTAINER_USER} ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" \
#     > /etc/sudoers.d/${CONTAINER_USER}-init-firewall \
#   && chmod 0440 /etc/sudoers.d/${CONTAINER_USER}-init-firewall

# Switch to container user
USER ${UID}:${GID}
WORKDIR ${WORKING_DIR}

# Copy project-specific configuration files
COPY --chown=${UID}:${GID} mise-config.toml ${MISEPATH}/config.toml
COPY --chown=${UID}:${GID} .cobra.yaml /home/${CONTAINER_USER}/.cobra.yaml

# Project-specific timezone (can be overridden per project)
ENV TZ="${TZ}"

# Define Mountpoint
# - COMMANDHISTORY: Persistent bash history across container restarts
# - MISEPATH: mise configuration, cache, and shims
# - PNPM_HOME: pnpm global store for faster package installation
# - WORKING_DIR: Project source code
# - node_modules: Node.js dependencies (improves performance on bind mounts)
# - CLAUDE_CODE_CONFIG: Claude Code settings and state
VOLUME [ "${COMMANDHISTORY}", "${MISEPATH}", "${PNPM_HOME}", "${WORKING_DIR}", "${WORKING_DIR}/node_modules", "${CLAUDE_CODE_CONFIG}" ]
