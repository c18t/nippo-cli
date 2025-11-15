# syntax=docker/dockerfile:1
# hadolint global ignore=DL3008,DL3046,DL3062

# ============================================
# Base Stage: Generic development environment
# ============================================
FROM golang:1.25.3-trixie AS base

# shell settings
SHELL ["/bin/bash", "-eo", "pipefail", "-c"]
ENV SHELL=bash
ENV TMP=/tmp
ENV TEMP=${TMP}

# go settings
ENV GOPATH="/go"

# workspace settings
ARG WORKSPACE_ROOT="/workspaces"
ARG WORKING_DIR="${WORKSPACE_ROOT}/app"

# set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

# container user settings
ARG CONTAINER_USER=user
ARG UID=1000
ARG GID=1000
ARG XDG_LOCAL="/home/${CONTAINER_USER}/.local"

# mise settings
ARG MISEPATH="/mise"
ENV MISE_DATA_DIR="${MISEPATH}"
ENV MISE_CONFIG_DIR="${MISEPATH}"
ENV MISE_CACHE_DIR="${MISEPATH}/cache"
ENV MISE_INSTALL_PATH="/usr/local/bin/mise"
ENV PATH="${MISEPATH}/shims:$PATH"

# pnpm settings
ENV PNPM_HOME="/home/${CONTAINER_USER}/.local/share/pnpm"

# claude code settings
ENV CLAUDE_CODE_CONFIG="/home/${CONTAINER_USER}/.claude"

# Install system packages, mise, create user, install Go tools
RUN <<EOF
# Add Debian sid repository to install the latest git version
echo 'deb http://deb.debian.org/debian sid main' > /etc/apt/sources.list.d/sid.list; \
printf 'Package: *\nPin: release a=sid\nPin-Priority: 100\n\n' \
  >  /etc/apt/preferences.d/git-from-sid
printf 'Package: git*\nPin: release a=sid\nPin-Priority: 990\n' \
  >> /etc/apt/preferences.d/git-from-sid
# Install mise dependencies and networking utilities for init-firewall
# cf. https://mise.jdx.dev/mise-cookbook/docker.html
#     https://code.claude.com/docs/en/devcontainer
apt-get update \
  && apt-get -y --no-install-recommends install \
    sudo curl git ca-certificates build-essential \
    less procps man-db unzip gnupg2 \
    iptables ipset iproute2 dnsutils aggregate jq \
  && rm -rf /var/lib/apt/lists/*
# Install mise (polyglot runtime manager) and allow mise self-update
curl https://mise.run | sh \
  && echo "%${CONTAINER_USER} ALL=(root) NOPASSWD: /usr/local/bin/mise self-update" > /etc/sudoers.d/${CONTAINER_USER}-mise-self-update \
  && chmod 0440 /etc/sudoers.d/${CONTAINER_USER}-mise-self-update
# Create non-root container user for security
groupadd -g ${GID} -o ${CONTAINER_USER}
useradd -m -s /bin/bash -u ${UID} -g ${GID} -o ${CONTAINER_USER}
# Create workspace, GOPATH, mise data, and configuration directories
mkdir -p ${GOPATH} ${MISEPATH} ${XDG_LOCAL} ${PNPM_HOME} ${WORKSPACE_ROOT} ${WORKING_DIR} ${WORKING_DIR}/node_modules ${CLAUDE_CODE_CONFIG}
# Grant ownership of all directories to the container user
chown -R ${CONTAINER_USER} ${GOPATH} ${MISEPATH} ${XDG_LOCAL} ${PNPM_HOME} ${WORKSPACE_ROOT} ${WORKING_DIR} ${WORKING_DIR}/node_modules ${CLAUDE_CODE_CONFIG}
# Configure persistent bash history across container restarts
SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history \
  && chown -R ${CONTAINER_USER} /commandhistory \
  && echo "${SNIPPET}" >> "/home/${CONTAINER_USER}/.bashrc"
# Install Go development tools
go install github.com/go-delve/delve/cmd/dlv@latest
go install golang.org/x/tools/gopls@latest
# Remove Go tool sources and prepare cache-volume mountpoint
rm -rf ${GOPATH}/pkg \
  && mkdir -p ${GOPATH}/pkg \
  && chown -R ${CONTAINER_USER} ${GOPATH}/pkg
EOF

# ============================================
# Project Stage: Project-specific configuration
# ============================================
FROM base AS project

# Re-declare ARGs needed in this stage
ARG CONTAINER_USER=user
ARG UID=1000
ARG GID=1000
ARG WORKSPACE_ROOT="/workspaces"
ARG WORKING_DIR="${WORKSPACE_ROOT}/app"
ARG MISEPATH="/mise"

# Project-specific timezone (can be overridden per project)
ARG TZ="Asia/Tokyo"
ENV TZ="${TZ}"

# Copy project-specific scripts
COPY start.sh post-create.sh init-firewall.sh /usr/local/bin/

# Configure scripts with appropriate permissions
RUN <<EOF
chmod +x /usr/local/bin/start.sh
chmod +x /usr/local/bin/post-create.sh
# skip firewall initialization
# chmod +x /usr/local/bin/init-firewall.sh \
#   && echo "%${CONTAINER_USER} ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/${CONTAINER_USER}-init-firewall \
#   && chmod 0440 /etc/sudoers.d/${CONTAINER_USER}-init-firewall
EOF

# Switch to container user
USER ${CONTAINER_USER}
WORKDIR ${WORKING_DIR}

# Copy project-specific configuration files
COPY --chown=${UID}:${GID} mise-config.toml ${MISEPATH}/config.toml
COPY --chown=${UID}:${GID} .cobra.yaml /home/${CONTAINER_USER}/.cobra.yaml

# ============================================
# Runtime configuration
# ============================================
# Delve debug server port
EXPOSE ${DEBUG_PORT_INTERNAL}

# Start delve DAP server
CMD ["start.sh"]
