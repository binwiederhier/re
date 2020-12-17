VERSION=0.1.0-alpha

.PHONY:

help:
	@echo "Build:"
	@echo "  make build             - Build"
	@echo "  make build-snapshot    - Build snapshot"
	@echo "  make build-simple      - Build (using go build, without goreleaser)"
	@echo "  make release           - Create a release"
	@echo "  make release-snapshot  - Create a test release"
	@echo "  make clean             - Clean build folder"

build: .PHONY
	goreleaser build --rm-dist

build-snapshot:
	goreleaser build --snapshot --rm-dist

build-simple: clean
	mkdir -p dist/re_linux_amd64
	go build \
		-o dist/re_linux_amd64/re \
		-ldflags \
		"-X main.version=${VERSION} -X main.commit=$(shell git rev-parse --short HEAD) -X main.date=$(shell date +%s)" \
		*.go

release:
	goreleaser release --rm-dist

release-snapshot:
	goreleaser release --snapshot --skip-publish --rm-dist

install:
	sudo rm -f /usr/bin/re
	sudo cp -a dist/re_linux_amd64/re /usr/bin/re

install-deb:
	sudo apt-get purge re || true
	sudo dpkg -i dist/*.deb

clean: .PHONY
	rm -rf dist
