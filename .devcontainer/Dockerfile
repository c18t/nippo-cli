# syntax=docker/dockerfile:1.4
FROM golang:1.22.5-bookworm

# install mise with apt
## xz-utils (for tar -xJv in mise install)
RUN <<EOF
apt-get update -y && apt-get install -y gpg sudo wget curl xz-utils less
gpg --keyserver hkps://keyserver.ubuntu.com --recv-keys 0x7413A06D
EOF

# create container user
ARG CONTAINER_USER=user
ARG UID=1000
ARG GID=1000
RUN <<EOF
groupadd -g ${GID} -o ${CONTAINER_USER}
useradd -m -s /bin/bash -u ${UID} -g ${GID} -o ${CONTAINER_USER}
EOF

# grant access to container user
COPY start.sh /usr/local/bin/start.sh
RUN <<EOF
chmod 755 /usr/local/bin/start.sh
chown -R ${CONTAINER_USER} /go
EOF

# configurations for container user
USER ${CONTAINER_USER}

ARG WORKING_DIR=/workspaces/app
WORKDIR ${WORKING_DIR}

# copy setting to home
COPY --chown=${UID}:${GID} mise-config.toml /home/${CONTAINER_USER}/.config/mise/config.toml
COPY --chown=${UID}:${GID} .cobra.yaml /home/${CONTAINER_USER}/.cobra.yaml

# go tools setting args
ARG DELVE_VERSION=latest

RUN <<EOF
: 'ensure the below is signed with the mise release key'
curl https://mise.jdx.dev/install.sh.sig | gpg --decrypt > install.sh
sh ./install.sh
: 'mise setup for bash'
echo 'export PATH=$PATH:~/.local/bin' >> ~/.bashrc
echo 'eval "$(mise activate bash)"' >> ~/.bashrc
: 'allow the execution of git commands in a container workspace'
git config --global safe.directory ${WORKING_DIR}
: 'install go tools'
go install github.com/go-delve/delve/cmd/dlv@${DELVE_VERSION}
EOF

# delve server port
ARG DEBUG_PORT=2345
EXPOSE ${DEBUG_PORT}

# start delve dap server
CMD ["start.sh"]
