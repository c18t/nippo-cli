# syntax=docker/dockerfile:1.4
FROM golang:1.21.5-bookworm

ARG DEBUG_PORT=2345
EXPOSE ${DEBUG_PORT}

ARG CONTAINER_USER=user
RUN useradd -m -s /bin/bash ${CONTAINER_USER}

COPY start.sh /usr/local/bin/start.sh
RUN <<EOF
chmod 755 /usr/local/bin/start.sh
chown -R ${CONTAINER_USER} /go
EOF

USER ${CONTAINER_USER}

ARG WORKING_DIR=/workspaces/go
WORKDIR ${WORKING_DIR}

RUN go install github.com/go-delve/delve/cmd/dlv@latest

CMD ["start.sh"]
