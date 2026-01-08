.PHONY: build clean install uninstall run test help

# 变量
APP_NAME := edge-client
BUILD_DIR := build
CMD_DIR := cmd/edge-client
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

# 默认目标
.DEFAULT_GOAL := help

## build: 编译 Edge Client
build:
	@echo "📦 编译 Edge Client..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) ./$(CMD_DIR)
	@echo "✅ 编译完成: $(BUILD_DIR)/$(APP_NAME)"

## build-all: 编译所有平台的二进制文件
build-all:
	@echo "📦 编译所有平台..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./$(CMD_DIR)
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 ./$(CMD_DIR)
	
	# Linux ARMv7
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-armv7 ./$(CMD_DIR)
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ./$(CMD_DIR)
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ./$(CMD_DIR)
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe ./$(CMD_DIR)
	
	@echo "✅ 所有平台编译完成"

## clean: 清理编译产物
clean:
	@echo "🧹 清理编译产物..."
	rm -rf $(BUILD_DIR)
	@echo "✅ 清理完成"

## install: 安装 Edge Client（需要 root 权限）
install: build
	@echo "📥 安装 Edge Client..."
	@if [ "$$(id -u)" -ne 0 ]; then \
		echo "⚠️  需要 root 权限，请使用: sudo make install"; \
		exit 1; \
	fi
	./install.sh

## uninstall: 卸载 Edge Client（需要 root 权限）
uninstall:
	@echo "🗑️  卸载 Edge Client..."
	@if [ "$$(id -u)" -ne 0 ]; then \
		echo "⚠️  需要 root 权限，请使用: sudo make uninstall"; \
		exit 1; \
	fi
	./uninstall.sh

## run: 运行 Edge Client（开发模式）
run: build
	@echo "🚀 运行 Edge Client..."
	./$(BUILD_DIR)/$(APP_NAME) start

## test: 运行测试
test:
	@echo "🧪 运行测试..."
	go test -v ./...

## fmt: 格式化代码
fmt:
	@echo "✨ 格式化代码..."
	go fmt ./...
	@echo "✅ 格式化完成"

## lint: 代码检查
lint:
	@echo "🔍 代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint 未安装，跳过检查"; \
		echo "安装: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

## deps: 下载依赖
deps:
	@echo "📦 下载依赖..."
	go mod download
	go mod tidy
	@echo "✅ 依赖下载完成"

## version: 显示版本信息
version:
	@echo "Edge Client"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"

## help: 显示帮助信息
help:
	@echo "PawStream Edge Client - Makefile"
	@echo ""
	@echo "可用命令:"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/  /'
	@echo ""
	@echo "示例:"
	@echo "  make build          # 编译"
	@echo "  make run            # 编译并运行"
	@echo "  sudo make install   # 安装到系统"
	@echo "  sudo make uninstall # 从系统卸载"
	@echo ""
