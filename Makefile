GO      := go
GOBUILD := $(GO) build
GOTEST  := $(GO) test

VTD_DIR := ./vtdxml
HELM_DIR = ./helm

TEST_MODULES = $(shell $(GO) list $(VTD_DIR)/...)

.PHONY: build
build:
	$(GOBUILD) ./...

.PHONY: lint
lint: golint helmlint

.PHONY: helmlint
helmlint:
	helm lint $(shell find $(HELM_DIR) -mindepth 1 -maxdepth 1 -type d)

.PHONY: golint
golint: vendor
	golangci-lint run

.PHONY: unit
unit:
	$(GOTEST) $(TEST_MODULES) \
		-cover \
		-coverprofile=c.out \
		-count=1 
	@cat c.out | \
		awk 'BEGIN {cov=0; stat=0;} $$3!="" { cov+=($$3==1?$$2:0); stat+=$$2; } \
		END {printf("Total coverage: %.2f%% of statements\n", (cov/stat)*100);}'
	go tool cover -html=c.out -o unit_test_coverage.html
