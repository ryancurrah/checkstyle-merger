current_dir = $(shell pwd)
version = $(shell printf '%s' $$(cat VERSION))

.PHONEY: lint
lint:
	golangci-lint run -v --enable-all --disable funlen,gochecknoglobals,lll ./...

.PHONEY: build
build:
	go build -o checkstyle-merger checkstyle-merger.go

.PHONEY: release
release:
	git tag ${version}
	git push --tags
	goreleaser --skip-validate --rm-dist

.PHONEY: clean
clean:
	rm -rf dist/
	rm -f checkstyle-merger

.PHONEY: install-tools-mac
install-tools-mac:
	brew install goreleaser/tap/goreleaser
