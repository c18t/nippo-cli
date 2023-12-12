#!/bin/bash

WORKING_DIR=${WORKING_DIR:-/workspaces/go}
CONTAINER_USER=${CONTAINER_USER:-user}
DEBUG_PORT=${DEBUG_PORT:-2345}

export HOME=/home/${CONTAINER_USER}

# カレントディレクトリの uid と gid を調べる
uid=${UID:-$(stat -c "%u" ${WORKING_DIR})}
gid=${GID:-$(stat -c "%g" ${WORKING_DIR})}

if [ "$uid" -ne 0 ]; then
  # ユーザーのuidをカレントディレクトリに合わせ、home, goディレクトリのuidを変更
  usermod -u $uid -o ${CONTAINER_USER} >/dev/null 2>&1
  chown -R $uid /go
fi

if [ "$gid" -ne 0 ]; then
  # ユーザーのgidをカレントディレクトリに合わせ、home, goディレクトリのgidを変更
  groupmod -g $gid -o ${CONTAINER_USER} >/dev/null 2>&1
  chgrp -R $gid ${HOME}
  chgrp -R $gid /go
fi

# uid/gid を指定して CMD 実行
exec setpriv --reuid=${CONTAINER_USER} --regid=${CONTAINER_USER} --init-groups /bin/bash -c "\
git config --global safe.directory ${WORKING_DIR} && \
id && $@"
