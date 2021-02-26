SHELL = bash
PROJECT_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_DIRTY := $(if $(shell git status --porcelain),+CHANGES)

GO_LDFLAGS ?= -X=github.com/seashell/agent/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)

CGO_ENABLED ?= 0

# User defined flags
STATIC := $(or $(STATIC),$(S)) ## If set to 1, build statically linked binary
DOCKER := $(or $(DOCKER),$(D)) ## If set to 1, build using docker container
OS := $(or $(OS),$(O)) # (coming soon) Define build target OS, e.g linux
ARCH := $(or $(ARCH),$(A)) # (coming soon) Define build target architecture, e.g amd64 

# Handle static builds
ifeq ($(STATIC),1)
	GO_LDFLAGS := "${GO_LDFLAGS} -linkmode external -extldflags -static"
endif

ifeq ($(DOCKER),1)
	CHECK_DOCKER := $(shell docker images --filter LABEL=com.seashell.builder=true -q)
ifeq ($(user),)
	# USER retrieved from env, UID from shell.
	HOST_USER ?= $(strip $(if $(USER),$(USER),nodummy))
	HOST_UID ?= $(strip $(if $(shell id -u),$(shell id -u),4000))
else
	# allow override by adding user= and/ or uid=  (lowercase!).
	# uid= defaults to 0 if user= set (i.e. root).
	HOST_USER = $(user)
	HOST_UID = $(strip $(if $(uid),$(uid),0))
endif
	BUILD_DOCKER := (docker build --label com.seashell.builder=true --build-arg HOST_UID=${HOST_UID} --build-arg HOST_USER=${HOST_USER} -t seashell_builder  . -f ./docker/Dockerfile.builder)
endif

# targets 
ALL_TARGETS += linux_amd64 \

default: help

build/linux_amd64/seashell: CMD='CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 \
								go build \
								-trimpath \
								-ldflags $(GO_LDFLAGS) \
								-o "$@"'							
build/linux_amd64/seashell: $(SOURCE_FILES) ## Build seashell for linux/amd64
	@echo "==> Building $@ with tags $(GO_TAGS)..."
ifeq ($(DOCKER),1)
ifeq ($(CHECK_DOCKER),)
	@echo "==> Building docker container builder image..."
	@$(call BUILD_DOCKER)
endif	
	docker run --rm -v ${PROJECT_ROOT}:${PROJECT_ROOT} --workdir=${PROJECT_ROOT} seashell_builder \
	/bin/sh -c ${CMD}
else
	@eval ${CMD}
endif

build/linux_arm64/seashell: # (coming soon) Build seashell for linux/arm64 
	@echo "==> Coming soon..."

build/linux_arm/seashell: # (coming soon) Build seashell for linux/arm
	@echo "==> Coming soon..."

.PHONY: ui
ui: CMD="go generate"
ui: ## Generate UI .go bindings
	@echo "==> Generating UI .go bindings..."
ifeq ($(DOCKER),1)
ifeq ($(CHECK_DOCKER),)
	@echo "==> Generating docker builder..."
	@$(call BUILD_DOCKER)
endif	
	docker run --rm -v ${PROJECT_ROOT}:${PROJECT_ROOT} --workdir=${PROJECT_ROOT} seashell_builder \
	/bin/sh -c ${CMD}
else
	@eval ${CMD}
endif

.PHONY: dev
dev: GOOS=$(shell go env GOOS)
dev: GOARCH=$(shell go env GOARCH)
dev: DEV_TARGET=build/$(GOOS)_$(GOARCH)/seashell
dev: ## Build for the current development platform
	@echo "==> Removing old development binary..."
	@rm -rf $(PROJECT_ROOT)/build
	@$(MAKE) --no-print-directory $(DEV_TARGET)

.PHONY: container
container: ## Build container with seashell binary inside
	@$(MAKE) ui dev STATIC=1 
	@echo "==> Building container image "seashell:latest" ..."
	@docker build -t seashell:latest . -f ./docker/Dockerfile.linux_amd64

.PHONY: release
release: clean ui $(foreach t,$(ALL_TARGETS),build/$(t)/seashell) ## Build all release packages which can be built on this platform
	@echo "==> Results:"
	@tree --dirsfirst $(PROJECT_ROOT)/build


.PHONY: clean
clean: ## Remove build artifacts
	@echo "==> Cleaning build artifacts..."
	@rm -rf "$(PROJECT_ROOT)/build/"
	@rm -rf "$(PROJECT_ROOT)/ui/build/"
	@rm -rf "$(PROJECT_ROOT)/ui/node_modules/"

HELP_FORMAT="    \033[36m%-25s\033[0m %s\n"
EG_FORMAT="    \033[36m%s\033[0m %s\n"
.PHONY: help
help: ## Display this usage information
	@echo "Valid targets:"
	@grep -E '^[^ ]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; \
			{printf $(HELP_FORMAT), $$1, $$2}'
	@echo ""
	@echo "This host will build the following targets if 'make release' is invoked:"
	@echo $(ALL_TARGETS) | sed 's/^/    /'
	@echo ""
	@echo "Valid flags:"
	@grep -E '^[^ ]+ :=.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = " :=.*?## "}; \
			{printf $(HELP_FORMAT), $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@printf $(EG_FORMAT) "~${PWD}" "$$ make ui dev DOCKER=1"
	@printf $(EG_FORMAT) "~${PWD}" "$$ make dev STATIC=1"
	@printf $(EG_FORMAT) "~${PWD}" "$$ make container DOCKER=1"
