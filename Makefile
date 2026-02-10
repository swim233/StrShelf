.PHONY: install build

build: install build_frontend build_backend

install:
	npm install

clean:
	@echo "cleaning dist files"
	rm -rv packages/web/dist
	rm -rv packages/api/dist

build_frontend:
	npm run build -w strshelf-web && \
	cp -rv packages/web/dist/ packages/api/

build_backend:
	cd packages/api && \
	mkdir -p ../../build/bin && \
	go build -v -o ../../build/bin/strshelf
