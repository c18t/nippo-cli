#!/bin/bash

# delveのdlv dapはクライアントのdisconnectでサーバープロセスが終了する制限事項があるため、ループで無限に起動させる
# cf. https://github.com/golang/vscode-go/blob/master/docs/debugging.md#remote-debugging
#     https://github.com/go-delve/delve/blob/master/Documentation/usage/dlv_dap.md
while :; do dlv dap -l 0.0.0.0:${DEBUG_PORT} --log --check-go-version=false; done
