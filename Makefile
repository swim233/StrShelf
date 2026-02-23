VERSION=$(shell git describe --tags --always --dirty)
GIT_COMMIT=$(shell git rev-parse HEAD | cut -c1-7)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=$(shell go version | awk '{print $$3}')
COMMIT_MESSAGE=$(shell git log -1 --pretty=%s | sed 's/"/\\"/g')
define _LDFLAGS
-ldflags "-X gopkg.ilharper.com/strshelf/api/lib.Version=$(VERSION) \
-X gopkg.ilharper.com/strshelf/api/lib.GitCommit=$(GIT_COMMIT) \
-X gopkg.ilharper.com/strshelf/api/lib.BuildTime=$(BUILD_TIME) \
-X gopkg.ilharper.com/strshelf/api/lib.GoVersion=$(GO_VERSION) \
-X 'gopkg.ilharper.com/strshelf/api/lib.CommitMessage=$(COMMIT_MESSAGE)' \
-X gopkg.ilharper.com/strshelf/api/lib.DebugModeStr=$1 \
$2"
endef

.PHONY: install build

build: install clean build_frontend build_backend
	@printf "\033[1;34mall target success complete!\033[0m\n"

install:
	@printf "\033[1;34minstalling npm module...\033[0m\n"
	pnpm install

clean:
	@printf "\033[1;34mcleaning dist files...\033[0m\n"
	rm -rfv packages/web/dist
	rm -rfv packages/api/dist

build_frontend:
	@printf "\033[1;34mpackaging frontend file...\033[0m\n"
	pnpm -C "packages/web" run build && \
	mkdir -p packages/api/dist/ && \
	cp -rv packages/web/dist/* packages/api/dist/

build_backend:
	@printf "\033[1;34mbuilding backend service...\033[0m\n"
	@printf "Building with DebugModeStr=false\n"
	cd packages/api && \
	mkdir -p ../../build/bin && \
	go build $(call _LDFLAGS,false,-s -w) -v -o ../../build/bin/strshelf

run: run_backend run_frontend

run_backend:
	@echo "Running with DebugModeStr=true"
	cd packages/api && \
	go run $(call _LDFLAGS,true,) -v main.go

run_frontend:
	cd packages/web && \
	pnpm -C "packages/web" run dev -- --host
