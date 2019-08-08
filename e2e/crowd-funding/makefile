SHELL := /bin/bash
NODE_PATH = $(shell ../scripts/find-node-or-install.sh)
PATH := $(NODE_PATH):$(shell echo $$PATH)

install:
	npm install

.PHONY: install	