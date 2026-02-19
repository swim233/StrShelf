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
-X gopkg.ilharper.com/strshelf/api/lib.DebugModeStr=$1"
endef

.PHONY: install build

build: install clean build_frontend build_backend
	@echo -e "\e[1;34mall target success complete!\e[0m"

install:
	@echo -e "\e[1;34minstalling npm module...\e[0m"
	npm install

clean:
	@echo -e "\e[1;34mcleaning dist files...\e[0m"
	rm -rfv packages/web/dist
	rm -rfv packages/api/dist

build_frontend:
	@echo -e "\e[1;34mpackaging frontend file...\e[0m"
	npm run build -w strshelf-web && \
	cp -rv packages/web/dist/ packages/api/

build_backend:
	@echo -e "\e[1;34mbuilding backend service...\e[0m"
	@echo "Building with DebugModeStr=false"
	cd packages/api && \
	mkdir -p ../../build/bin && \
	go build $(call _LDFLAGS,false) -v -o ../../build/bin/strshelf

run: run_backend run_frontend

run_backend:
	@echo "Running with DebugModeStr=true"
	cd packages/api && \
	go run $(call _LDFLAGS,true) -v main.go

run_frontend:
	cd packages/web && \
	npm run dev -- --host
