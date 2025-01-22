SHELL:=bash
.PHONY: help start stop logs
UID := $(shell id -u)
GID := $(shell id -g)
export UID
export GID

white := \033[0;1m
green := \033[0;32m
red := \033[0;31m
light_red := \033[1;31m
yellow := \033[1;33m
cyan := \033[0;36m
nc := \033[0m	

help: ## Show this help with the list of commands
	@echo -e "\n\033[0;1m▼▼▼ PROFANITY AI CLEAN CHAT ▼▼▼\033[0m"
	@echo -e "\nYou can override some variables by editing this Makefile or by using \"make <target> <var>=<value>\"."
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[0;33m%s\033[0m\n", substr($$0, 5) } \
		/^##=/ { printf "\n\033[1;35m▼ %s ▼\033[0m\n", substr($$0, 5) } \
		/^##n/ { printf "\n" } \
		' $(MAKEFILE_LIST)


define init_project
	@echo -e "\n${cyan}🚀 Initializing project...${nc}"
	@docker compose build --no-cache
	@docker compose up -d
	@echo -e "\n${green}✅ Project initialized!${nc}\n"
endef

define start_project
	@echo -e "\n${cyan}🚀 Starting project...${nc}"
	@docker compose up -d
	@echo -e "\n${green}✅ Project started!${nc}\n"
endef

define stop_project
	@echo -e "\n${cyan}🚀 Stopping project...${nc}"
	@docker compose down
	@echo -e "\n${green}✅ Project stopped!${nc}\n"
endef

define logs_project
	@echo -e "\n${cyan}🚀 Showing logs...${nc}"
	@docker compose logs -f
endef

define deno_dev
	@echo -e "\n${cyan}🚀 Starting development server...${nc}"
	@docker compose exec frontend deno run dev
endef

define deno_install
	@echo -e "\n${cyan}🚀 Installing dependencies...${nc}"
	@docker compose exec frontend deno install --allow-scripts=npm:@sveltejs/kit@2.15.2
endef

init: ## Initialize the project
	$(call init_project)

start: ## Start the project
	$(call start_project)

stop: ## Stop the project
	$(call stop_project)

logs: ## Show the logs
	$(call logs_project)

dev: ## Start the development server
	$(call deno_dev)

install: ## Install dependencies
	$(call deno_install)

clean: ## Clean the Project
	@echo -e "\n${cyan}🚀 Cleaning project...${nc}"
	@echo -r "Are you sure you want to clean the project? [y/N] " && read ans && [ $${ans:-N} = y ]
	@echo -e "\n${cyan}Deleting dockers...${nc}"
	@docker compose down -v
	@echo -e "\n${cyan}Removing frontend dependencies...${nc}"
	@rm -rf frontend/node_modules frontend/.svelte-kit
	@echo -e "\n${cyan}Removing backend dependencies...${nc}"
	@echo -e "\n${green}✅ Project cleaned!${nc}\n"
