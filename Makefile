BUILD_COMMIT=$(shell git log --format="%H" -n 1)
LDFLAGS="-X 'main.BuildGitHash=$(BUILD_COMMIT)'"

HAS_LINT := $(shell command -v golangci-lint;)
HAS_IMPORTS := $(shell command -v goimports;)

PROJECT = github.com/t0pep0/GB_best_go1
GO_PKG = $(shell go list $(PROJECT)/hw1/...)

bootstrap:
ifndef HAS_LINT
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0
endif
ifndef HAS_IMPORTS
	go install golang.org/x/tools/cmd/goimports
endif

build: test
	go build -ldflags $(LDFLAGS) -o bin/bot ./main.go

test:
	@echo "+ $@"
	@go list -f '"go test -v {{.Dir}}"' $(GO_PKG) | xargs -L 1 sh -c

fmt:
	@echo "+ $@"
	@go list -f '"gofmt -w -s -l {{.Dir}}"' $(GO_PKG) | xargs -L 1 sh -c

imports:
	@echo "+ $@"
	@go list -f '"goimports -w {{.Dir}}"' ${GO_PKG} | xargs -L 1 sh -c

lint: bootstrap
	@echo "+ $@"
	@golangci-lint run ./...

.PHONY: bootstrap \
	build \
	test \
	fmt \
	imports \
	lint