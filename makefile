# 默认变量定义
BINARY=bm

# 编译默认目标
all: build

# 构建可执行文件
build:
	GOOS=linux
	GOARCH=amd64
	go build -o ${BINARY} -ldflags="-w -s" -trimpath cmd/main.go

# 使用 swag 生成 swagger 文档
swagger:
	swag init -g cmd/main.go

wire:
	go run cmd/gen_wire/main.go
	go generate ./internal/...

# 清理构建产物
clean:
	rm -f ${BINARY}

# 运行程序
run:
	go run cmd/main.go serve

# 帮助信息
help:
	@echo "Usage:"
	@echo "  make build      - 编译程序"
	@echo "  make swagger    - 生成 swagger 文档 (需要安装 swag)"
	@echo "  make run        - 运行程序"
	@echo "  make clean      - 清理编译产物"
	@echo "  make wire       - 更新wire依赖，并重新生成wire_gen"