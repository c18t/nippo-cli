# syntax=docker/dockerfile:1.4
FROM golang:1.25.3-trixie

# shell settings
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# workspace settings
ARG WORKING_DIR=/workspaces/app

# set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

# mise settings
ENV MISE_DATA_DIR="/mise"
ENV MISE_CONFIG_DIR="/mise"
ENV MISE_CACHE_DIR="/mise/cache"
ENV MISE_INSTALL_PATH="/usr/local/bin/mise"
ENV PATH="/mise/shims:$PATH"
# ENV MISE_VERSION="..."

# claude code settings
ARG TZ=
ENV TZ="${TZ}"
ARG CLAUDE_CODE_VERSION=latest

# container user settings
ARG CONTAINER_USER=user
ARG UID=1000
ARG GID=1000

# Copy scripts
COPY init-firewall.sh /usr/local/bin/
COPY start.sh /usr/local/bin/

RUN <<EOF
: 'install mise deendencies and utils for init-firewall'
apt-get update \
  && apt-get -y --no-install-recommends install \
    sudo curl git ca-certificates build-essential \
    less procps fzf man-db unzip gnupg2 \
    iptables ipset iproute2 dnsutils aggregate jq \
  && rm -rf /var/lib/apt/lists/*
: 'install mise'
curl https://mise.run | sh
: 'create container user'
groupadd -g ${GID} -o ${CONTAINER_USER}
useradd -m -s /bin/bash -u ${UID} -g ${GID} -o ${CONTAINER_USER}
: 'create workspace and config directories and set permissions'
mkdir -p ${WORKING_DIR} /mise /home/${CONTAINER_USER}/.claude \
  && chown -R ${CONTAINER_USER} ${WORKING_DIR} /mise /home/${CONTAINER_USER}/.claude
: 'grant access to container user'
chmod 755 /usr/local/bin/start.sh
chown -R ${CONTAINER_USER} /go
: 'set up init-firewall script'
chmod +x /usr/local/bin/init-firewall.sh \
  && echo "%${CONTAINER_USER} ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/${CONTAINER_USER}-init-firewall \
  && chmod 0440 /etc/sudoers.d/${CONTAINER_USER}-init-firewall
: 'grant chown execution to container user'
echo "%${CONTAINER_USER} ALL=(root) NOPASSWD: /usr/bin/chown" > /etc/sudoers.d/${CONTAINER_USER}-chown \
  && chmod 0440 /etc/sudoers.d/${CONTAINER_USER}-chown
: 'persist bash history'
SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history \
  && chown -R ${CONTAINER_USER} /commandhistory
EOF

# configurations for container user
USER ${CONTAINER_USER}
WORKDIR ${WORKING_DIR}

# copy setting to home
COPY --chown=${UID}:${GID} mise-config.toml ${MISE_CONFIG_DIR}/config.toml
COPY --chown=${UID}:${GID} .cobra.yaml /home/${CONTAINER_USER}/.cobra.yaml

# go tools setting args
ARG DELVE_VERSION=latest
ARG GOPLS_VERSION=latest

RUN <<EOF
: 'allow the execution of git commands in a container workspace'
git config --global safe.directory ${WORKING_DIR}
: 'install go tools'
go install github.com/go-delve/delve/cmd/dlv@${DELVE_VERSION}
go install golang.org/x/tools/gopls@${GOPLS_VERSION}
EOF

# delve server port
ARG DEBUG_PORT=2345
EXPOSE ${DEBUG_PORT}
# start delve dap server
CMD ["start.sh"]
