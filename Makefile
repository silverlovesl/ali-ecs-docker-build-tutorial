# =============================================================================
#  ali-ecs-docker-build-tutorial — monorepo task runner
# =============================================================================
#  Make targets wrap the root `package.json` scripts (which in turn drive
#  Turborepo, oxc and Go). Run `make` to list every target.
#
#  Scope Turborepo targets with FILTER, for example:
#      make build FILTER=@repo/web
#      make test  FILTER=@repo/api
# =============================================================================

.DEFAULT_GOAL := help
SHELL := /bin/bash

BUN   ?= bun
GO    ?= go
TURBO ?= $(BUN)x turbo

# Optional Turborepo `--filter` value, e.g. `make lint FILTER=@repo/web`
FILTER ?=
TURBO_FILTER := $(if $(FILTER),--filter=$(FILTER),)

WEB := apps/web
API := apps/api

C_RESET := \033[0m
C_BOLD  := \033[1m
C_CYAN  := \033[36m

.PHONY: help install hooks setup \
        dev dev-web dev-api \
        build build-web build-api preview start-api generate \
        lint lint-fix format format-check typecheck test test-api check \
        docker-build docker-up docker-down docker-logs \
        clean clean-all versions

##@ Setup

install: ## Install all workspace dependencies (also installs git hooks)
	$(BUN) install

hooks: ## Install or refresh the Lefthook git hooks
	$(BUN)x lefthook install

setup: install hooks ## Install dependencies and git hooks

##@ Development

dev: ## Start web + api with hot reload
	$(TURBO) run dev $(TURBO_FILTER)

dev-web: ## Start only the web dev server (http://localhost:3000)
	$(TURBO) run dev --filter=@repo/web

dev-api: ## Start only the API with air hot reload (http://localhost:8080)
	$(TURBO) run dev --filter=@repo/api

##@ Build and run

build: ## Build every workspace
	$(TURBO) run build $(TURBO_FILTER)

build-web: ## Build the web app into apps/web/dist
	$(TURBO) run build --filter=@repo/web

build-api: ## Build the API binary into apps/api/bin/api
	$(TURBO) run build --filter=@repo/api

preview: ## Serve the production web build locally
	cd $(WEB) && $(BUN) run preview

start-api: build-api ## Build and run the compiled API binary
	cd $(API) && $(BUN) run start

generate: ## Regenerate the TanStack Router route tree
	cd $(WEB) && $(BUN)x --bun tsr generate

##@ Docker

docker-build: ## Build both images locally (ali-ecs-docker-build-tutorial-web, -api)
	docker compose build

docker-up: ## Build if needed and start both containers (http://localhost:3001)
	docker compose up -d --build

docker-down: ## Stop and remove the containers
	docker compose down

docker-logs: ## Follow container logs
	docker compose logs -f

##@ Quality

lint: ## Lint TS/TSX with oxlint and Go with `go vet`
	$(TURBO) run lint $(TURBO_FILTER)

lint-fix: ## Apply oxlint fixes
	$(TURBO) run lint:fix $(TURBO_FILTER)

format: ## Format the repo with oxfmt and gofmt
	$(BUN) run format

format-check: ## Check formatting without writing (CI friendly)
	$(BUN) run format:check

typecheck: ## Type-check TS with tsc and Go with `go build`
	$(TURBO) run typecheck $(TURBO_FILTER)

test: ## Run every test suite
	$(TURBO) run test $(TURBO_FILTER)

test-api: ## Run only the Go tests
	$(TURBO) run test --filter=@repo/api

check: format-check lint typecheck test ## Run everything CI runs (format, lint, types, tests)

##@ Cleanup

clean: ## Remove build output and generated files
	$(TURBO) run clean

clean-all: clean ## Also remove node_modules
	rm -rf node_modules
	rm -rf .build
	rm -rf .turbo

##@ Misc

versions: ## Print the versions of the toolchain
	@printf "bun    %s\n" "$$($(BUN) --version)"
	@printf "go     %s\n" "$$($(GO) version | awk '{print $$3}')"
	@printf "turbo  %s\n" "$$($(TURBO) --version)"

##@ Help

help: ## Show this help
	@printf "\n$(C_BOLD)ali-ecs-docker-build-tutorial$(C_RESET)\n\n"
	@printf "Usage: make $(C_CYAN)<target>$(C_RESET) [FILTER=<package>]\n"
	@awk 'BEGIN {FS = ":.*?## "} \
	  /^##@/ { printf "\n$(C_BOLD)%s$(C_RESET)\n", substr($$0, 5) } \
	  /^[a-zA-Z0-9_-]+:.*?## / { printf "  $(C_CYAN)%-16s$(C_RESET) %s\n", $$1, $$2 }' \
	  $(MAKEFILE_LIST)
	@printf "\nExamples:\n"
	@printf "  make dev\n"
	@printf "  make test FILTER=@repo/api\n"
	@printf "  make build-web && make start-api\n\n"
