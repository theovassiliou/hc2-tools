# Inspired by github.com/influxdata/telegraf
ifeq ($(OS), Windows_NT)
	VERSION := $(shell git describe --exact-match --tags 2>nil)
	HOME := $(HOMEPATH)
	CGO_ENABLED ?= 0
	export CGO_ENABLED
else
	VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
endif

PREFIX := /usr/local
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git rev-parse --short HEAD)
GOFILES ?= $(shell git ls-files '*.go')
GOFMT ?= $(shell gofmt -l -s $(GOFILES))
BUILDFLAGS ?=

ifdef GOBIN
PATH := $(GOBIN):$(PATH)
else
PATH := $(subst :,/bin:,$(shell go env GOPATH))/bin:$(PATH)
endif

LDFLAGS := $(LDFLAGS) -X main.commit=$(COMMIT) -X main.branch=$(BRANCH)
ifdef VERSION
	LDFLAGS += -X main.version=$(VERSION)
endif

.PHONY: all
all:
	@$(MAKE) --no-print-directory deps
	@$(MAKE) --no-print-directory hc2

.PHONY: deps
deps:
	go mod vendor

.PHONY: hc2

hc2:
	go build -ldflags "$(LDFLAGS)" ./cmd/hc2UploadScene
	go build -ldflags "$(LDFLAGS)" ./cmd/hc2DownloadScene
	go build -ldflags "$(LDFLAGS)" ./cmd/hc2SceneInteract


.PHONY: go-install
go-install:
	go install -ldflags "-w -s $(LDFLAGS)" ./cmd/hc2UploadScene
	go install -ldflags "-w -s $(LDFLAGS)" ./cmd/hc2DownloadScene
	go install -ldflags "-w -s $(LDFLAGS)" ./cmd/hc2SceneInteract

.PHONY: install
install: hc2
	mkdir -p $(DESTDIR)$(PREFIX)/bin/
	cp hc2UploadScene $(DESTDIR)$(PREFIX)/bin/
	cp hc2DownloadScene $(DESTDIR)$(PREFIX)/bin/
	cp hc2SceneInteract $(DESTDIR)$(PREFIX)/bin/

.PHONY: test
test:
	go test -short ./...

.PHONY: fmt
fmt:
	@gofmt -s -w $(GOFILES)

.PHONY: fmtcheck
fmtcheck:
	@if [ ! -z "$(GOFMT)" ]; then \
		echo "[ERROR] gofmt has found errors in the following files:"  ; \
		echo "$(GOFMT)" ; \
		echo "" ;\
		echo "Run make fmt to fix them." ; \
		exit 1 ;\
	fi

.PHONY: vet
vet:
	@echo 'go vet $$(go list ./...)'
	@go vet $$(go list ./...) ; if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "go vet has found suspicious constructs. Please remediate any reported errors"; \
		echo "to fix them before submitting code for review."; \
		exit 1; \
	fi

.PHONY: check
check: fmtcheck vet

.PHONY: test-all
test-all: fmtcheck vet
	go test ./...

.PHONY: package
package:
	./scripts/build.py --package --platform=all --arch=all

.PHONY: package-release
package-release:
	./scripts/build.py --release --package --platform=all --arch=all

.PHONY: package-nightly
package-nightly:
	./scripts/build.py --nightly --package --platform=all --arch=all 

.PHONY: clean
clean:
	rm -f $(GOPATH)/bin/hc2UploadScene
	rm -f $(GOPATH)/bin/hc2UploadScene.exe
	rm -f $(GOPATH)/bin/hc2DownloadScene
	rm -f $(GOPATH)/bin/hc2DownloadScene.exe
	rm -f $(GOPATH)/bin/hc2SceneInteract
	rm -f $(GOPATH)/bin/hc2SceneInteract.exe
	rm -f ./hc2UploadScene
	rm -f ./hc2DownloadScene
	rm -f ./hc2SceneInteract

.PHONY: docker-image

.PHONY: static
static:
	@echo "Building static linux binary hc2UploadScene"
	@CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" ./cmd/hc2UploadScene

	@echo "Building static linux binary hc2DownloadScene"
	@CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" ./cmd/hc2DownloadScene

	@echo "Building static linux binary hc2SceneInteract"
	@CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" ./cmd/hc2SceneInteract

