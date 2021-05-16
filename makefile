# メタ情報

export GO111MODULE=on
# 開発に必要な依存をインストールする
## Setup for Development env
.PHONY: setup-tools
setup-tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.37.0
	go install github.com/Songmu/make2help/cmd/make2help@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/cespare/reflex@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install golang.org/x/tools/cmd/stringer@latest
	go install github.com/moznion/gonstructor/cmd/gonstructor@latest


## run server with hot-reload
dev:
	reflex -r '\.go' -s go run main.go

.PHONY: build
build:
	go build main.go


