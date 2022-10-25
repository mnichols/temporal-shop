.PHONY: clean test bins lint tools
PROJECT_ROOT = github.com/temporalio/temporal-shop

# default target
default: clean test bins

TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(TAG:v%=%)
APPPKG := $(PROJECT_ROOT)/app
LINKER_FLAGS := -X $(APPPKG).BuildDate=$(DATE) -X $(APPPKG).Commit=$(COMMIT) -X $(APPPKG).Version=$(VERSION)

web:
	@go build -ldflags "$(LINKER_FLAGS)" -o web ./web/bff/cmd/bff/*.go

bins: web

test:
	go test -race -timeout=5m -cover -count=1  ./...

clean:
	@rm -rf ./web

define build
	@echo "building release for $(1) $(2) $(3)..."
	@mkdir -p releases
	@GOOS=$(2) GOARCH=$(3) go build -ldflags "-w $(LINKER_FLAGS)" -o releases/$(1)_$(2)_$(3)$(4) ./cmd/web/*.go
	@tar -cvzf releases/$(1)_$(2)_$(3).tar.gz releases/$(1)_$(2)_$(3)$(4) &>/dev/null
endef

tools:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2
	@GO111MODULE=off go get -u github.com/golang/mock/mockgen

lint:
	golangci-lint run