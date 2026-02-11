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
	@echo -e "\e[1;34mbuilding go binary file...\e[0m"
	cd packages/api && \
	mkdir -p ../../build/bin && \
	go build -ldflags="-X 'gopkg.ilharper.com/strshelf/api/config.DebugModeStr=false'" -v -o ../../build/bin/strshelf

run: run_backend run_frontend

run_backend:
	cd packages/api && \
	go run -ldflags="-X 'gopkg.ilharper.com/strshelf/api/config.DebugModeStr=true'" main.go

run_frontend:
	cd packages/web && \
	npm run dev -- --host
