# 出力先のディレクトリ
BINDIR:=bin

# バージョン
VER_TAG:=$(shell git describe --tag)
VER_REV:=$(shell git rev-parse --short HEAD)
VERSION:=${VER_TAG}+${VER_REV}

# ldflag
GO_LDFLAGS_NOST:=-s -w
GO_LDFLAGS_STATIC:=-extldflags "-static"
GO_LDFLAGS_VERSION:=-X "main.version=${VERSION}"
GO_LDFLAGS:=$(GO_LDFLAGS_NOST) $(GO_LDFLAGS_STATIC) $(GO_LDFLAGS_VERSION)

# go build
GO_BUILD:=-ldflags '${GO_LDFLAGS}' -trimpath

# ビルドタスク
all: build

.PHONY: build
build: app

.PHONY: debug
debug: BINDIR:=/tmp
debug: GO_BUILD:=-gcflags='all=-N -l' -ldflags '${GO_LDFLAGS_VERSION}'
debug: app

.PHONY: app
app: nippo/nippo.go
	go build -o $(BINDIR)/$(<F:.go=) $(GO_BUILD) $^
