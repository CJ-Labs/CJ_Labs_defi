# Git 相关变量
# 定义变量 GITCOMMIT，其值为当前 Git 仓库的最新提交哈希（commit hash）
GITCOMMIT := $(shell git rev-parse HEAD)
# 定义变量 GITDATE，其值为最新提交的 Unix 时间戳。
GITDATE := $(shell git show -s --format='%ct')

# 链接标志
# 将 -X main.GitCommit=<commit_hash> 添加到 LDFLAGSSTRING 变量。
LDFLAGSSTRING += -X main.GitCommit=$(GITCOMMIT)
# 将 -X main.GitDate=<timestamp> 追加到 LDFLAGSSTRING。
LDFLAGSSTRING += -X main.GitDate=$(GITDATE)
# 定义 LDFLAGS 变量，其值为 -ldflags "<LDFLAGSSTRING>"。
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

# ABI 文件路径和输出目录
# 定义变量 EVENT_ABI_ARTIFACTS，表示输入 ABI 文件的通配符路径。
EVENT_ABI_ARTIFACTS := ./abis/arbitrum/*.json
# 定义变量 OUTPUT_DIR，表示生成文件的输出目录。
OUTPUT_DIR := ./event/arbitrum

# 确保输出目录存在
$(shell mkdir -p $(OUTPUT_DIR))

# 默认目标
all: binding-event

# 批量处理 ABI 文件
binding-event:
	@for abi_file in $(EVENT_ABI_ARTIFACTS); do \
		base_name=$$(basename "$$abi_file" .json); \
		type_name="$${base_name}EventManager"; \
		pkg_name=$$(echo "$$base_name" | tr '[:upper:]' '[:lower:]'); \
		output_file="$(OUTPUT_DIR)/$${pkg_name}/$${base_name}.go"; \
		echo "Processing $$abi_file -> $$output_file"; \
		mkdir -p "$(OUTPUT_DIR)/$${pkg_name}"; \
		temp_file=$$(mktemp); \
		cat "$$abi_file" | jq .abi > "$$temp_file"; \
		abigen --pkg "$$pkg_name" \
			--abi "$$temp_file" \
			--out "$$output_file" \
			--type "$$type_name"; \
		rm "$$temp_file"; \
	done
# 清理生成的文件
clean:
	rm -f $(OUTPUT_DIR)/*.go

# 伪目标
.PHONY: \
	all \
	clean