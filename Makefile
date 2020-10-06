.ONESHELL:
.DEFAULT_GOAL := list
VERSION := 5
ROOT_DIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
DEPLOY_DIR := $(HOME)/.terraform.d/plugins/onemoregame.com/terraform/discord/0.0.$(VERSION)/linux_amd64


# https://stackoverflow.com/a/26339924/11547115
.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

build: $(ROOT_DIR)/terraform-provider-discord

deploy: build $(DEPLOY_DIR)/terraform-provider-discord

$(ROOT_DIR)/terraform-provider-discord: *.go discord/*.go go.mod go.sum Makefile
	cd "$(ROOT_DIR)"
	go build -o terraform-provider-discord

$(DEPLOY_DIR)/terraform-provider-discord: $(ROOT_DIR)/terraform-provider-discord
	mkdir -p "$(DEPLOY_DIR)"
	cp "$(ROOT_DIR)/terraform-provider-discord" "$(DEPLOY_DIR)"
