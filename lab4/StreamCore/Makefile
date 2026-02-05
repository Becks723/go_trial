
# 环境变量
DIR=$(shell pwd)
CMD_PATH=$(DIR)/cmd
CONFIG_PATH=$(DIR)/config
IDL_PATH=$(DIR)/idl
OUTPUT_PATH=$(DIR)/output
API_PATH=$(DIR)/cmd/api

# 服务名
SERVICES=api user video interaction social
service=$(word 1, $@)

PREFIX="[MAKEFILE]"

# 构建指定的服务并运行，如果仅构建需要加上 BUILD_ONLY 参数
.PHONY: $(SERVICES)
$(SERVICES):
	  @echo "$(PREFIX) Building $(service)..."
	  @mkdir -p output
	  @bash $(DIR)/docker/script/build.sh $(service)
	  @echo "$(PREFIX) Build $(service) completed"
ifndef BUILD_ONLY
	  @echo "$(PREFIX) Starting $(service)..."
	  @export SERVICE=$(service); bash ./docker/script/entrypoint.sh
endif

# 启动必要的环境，如 mysql和etcd 等
.PHONY: env-up
env-up:
	@docker compose -f ./docker/docker-compose.yml up -d

# 关闭环境
.PHONY: env-down
env-down:
	@cd ./docker && docker compose down


