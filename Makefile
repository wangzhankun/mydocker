# Makefile 基础教程
# https://www.kancloud.cn/kancloud/make-command/45596
# https://seisman.github.io/how-to-write-makefile/overview.html

.PHONY:build
build:
	go build -v .

.PHONY:run
run:
	go run main.go

.PHONY:test
test:
	go test ./...

# 执行检测 指定超时时间和配置文件
.PHONY:lint
lint:
	golangci-lint run --timeout=5m --config ./.golangci.yaml
# 清理缓存
.PHONY:lint-clean
lint-clean:
	golangci-lint cache clean

