#! /usr/bin/make
#(C) Copyright 2021 Hewlett Packard Enterprise Development LP
# Inspiration from https://github.com/rightscale/go-boilerplate/blob/master/Makefile

NAME=$(shell find cmd -name ".gitkeep_provider" -exec dirname {} \; | sort -u | sed -e 's|cmd/||')
VERSION=0.0.1
# Change DUMMY_PROVIDER below to reflect the name of the service under development.  The
# value of this variable is used in LOCAL_LOCATION, and is also used in the
DUMMY_PROVIDER=vmaas
LOCAL_LOCATION=~/.local/share/terraform/plugins/terraform.example.com/$(DUMMY_PROVIDER)/hpegl/$(VERSION)/linux_amd64/

# Stuff that needs to be installed globally (not in vendor)
DEPEND=

# Will get the branch name
SYMBOLIC_REF=$(shell if [ -n "$$CIRCLE_TAG" ] ; then echo $$CIRCLE_TAG; else git symbolic-ref HEAD | cut -d"/" -f 3; fi)
COMMIT_ID=$(shell git rev-parse --verify HEAD)
DATE=$(shell date +"%F %T")

PACKAGE := $(shell git remote get-url origin | sed -e 's|http://||' -e 's|^.*@||' -e 's|.git||' -e 's|:|/|')
VERSION_PACKAGE=$(PACKAGE)/pkg/cmd/$@
VFLAG=-X '$(VERSION_PACKAGE).name=$@' \
      -X '$(VERSION_PACKAGE).version=$(SYMBOLIC_REF)' \
      -X '$(VERSION_PACKAGE).buildDate=$(DATE)' \
      -X '$(VERSION_PACKAGE).buildSha=$(COMMIT_ID)'
TAGS=

# kelog issue: https://github.com/rjeczalik/notify/issues/108
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	TAGS=-tags kqueue
endif
TMPFILE := $(shell mktemp)

LOCALIZATION_FILES := $(shell find . -name \*.toml | grep -v vendor | grep -v ./bin)

$(NAME): $(shell find . -name \*.go)
	CGO_ENABLED=0 go build $(TAGS) -ldflags "$(VFLAG)" -o build/$@ .

default: all
.PHONY: default

generate:
	go generate ./...

vendor: generate go.mod go.sum
	go mod download

update up: really-clean vendor
.PHONY: update up

clean:
	rm -rf gathered_logs build .vendor/pkg $(testreport_dir) $(coverage_dir)
.PHONY: clean

really-clean clean-all cleanall: clean
	rm -rf vendor
.PHONY: really-clean clean-all cleanall

procs := $(shell grep -c ^processor /proc/cpuinfo 2>/dev/null || echo 1)
# TODO make --debug an option

lint: vendor golangci-lint-config.yaml
	@golangci-lint --version
	golangci-lint run --config golangci-lint-config.yaml
.PHONY: lint

testreport_dir := test-reports
unit-test: generate
	@go test `go list ./... | grep -v github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/acceptance_test`
.PHONY: test

coverage_dir := coverage/go
coverage: vendor
	@mkdir -p $(coverage_dir)/html
	go test -coverpkg=./... -coverprofile=$(coverage_dir)/coverage.out -v $$(go list ./... | grep -v /vendor/)
	@go tool cover -html=$(coverage_dir)/coverage.out -o $(coverage_dir)/html/main.html;
	@echo "Generated $(coverage_dir)/html/main.html";
.PHONY: coverage

ACC_TEST_FILE_LOCATION=github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/acceptance_test
acceptance:
	@if [ "${case}" != "" ]; then \
		TF_ACC=true go test -run $(case) -v -timeout=2000s -cover $(ACC_TEST_FILE_LOCATION); \
	else \
		TF_ACC=true go test -v -timeout=9000s -cover $(ACC_TEST_FILE_LOCATION); \
	fi

build: vendor $(NAME)
.PHONY: build

install: build $(NAME)
	# terraform >= v0.13
	mkdir -p $(LOCAL_LOCATION)
	cp build/$(NAME) $(LOCAL_LOCATION)
.PHONY: install

v := latest
sdk:
	rm -rf vendor
	@go get github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk@$v
	go mod vendor
.PHONY: v

tflint:
	@terraform fmt -recursive ./examples/

all: lint test
.PHONY: all

tools:
	go env -w GO111MODULE=on
	go env -w GOPRIVATE="github.com/hpe-hcss/*"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
.PHONY: tools
