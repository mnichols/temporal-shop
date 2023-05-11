.PHONY: clean test bins lint tools
GITHUB_REPOSITORY = temporalio/temporal-shop

# default target
default: clean test bins

TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
COMMIT := $(shell git rev-parse --short HEAD )
SHORT_SHA := $(echo $COMMIT | cut -c 1-8)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(TAG:v%=%)
BFF_BUILD_PKG := github.com/$(GITHUB_REPOSITORY)/web/bff/build
LINKER_FLAGS := "-X '${BFF_BUILD_PKG}.Commit=${COMMIT}' -X '${BFF_BUILD_PKG}.Version=${VERSION}' -X '${BFF_BUILD_PKG}.BuildDate=${DATE}'"

out:
	mkdir -p out

bff: out
	@go work sync;cd web; go build -ldflags ${LINKER_FLAGS} -o ../out/bff ./bff/cmd/bff/main.go

domain: out
	@cd services/go; go build -ldflags ${LINKER_FLAGS} -o ../out/domain ./cmd/temporal_shop/main.go

bins: bff domain

genapi:
	@cd proto; buf generate

gengql:
	@cd web/ui; npm run codegen
	@cd web/bff/internal/gql; go run github.com/99designs/gqlgen --config gqlgen.yml generate
	#@cd services/go/internal/pubsub; go run github.com/Khan/genqlient ./genqlient.yaml

test:
	go test -race -timeout=5m -cover -count=1  ./...

clean:
	@rm -rf out

lint:
	golangci-lint run ./web

login:
	JWT = $(http --verify=no https://localhost:8080/api/login email=mike.nichols@temporal.io | jq -r .token)
	echo ${JWT}
