# syntax=docker/dockerfile:1.4
FROM golang:1.21.5-bookworm

ARG DEBUG_PORT=2345
EXPOSE ${DEBUG_PORT}

ARG CONTAINER_USER=user
ARG UID=1000
ARG GID=1000
RUN <<EOF
groupadd -g ${GID} -o ${CONTAINER_USER}
useradd -m -s /bin/bash -u ${UID} -g ${GID} -o ${CONTAINER_USER}
EOF

COPY start.sh /usr/local/bin/start.sh
RUN <<EOF
chmod 755 /usr/local/bin/start.sh
chown -R ${CONTAINER_USER} /go
EOF

USER ${CONTAINER_USER}

ARG WORKING_DIR=/workspaces/go
WORKDIR ${WORKING_DIR}

RUN <<EOF
git config --global safe.directory ${WORKING_DIR}
go install github.com/go-delve/delve/cmd/dlv@latest
EOF

CMD ["start.sh"]
