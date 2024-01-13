.SILENT :

export GO111MODULE=on

# App name
APPNAME=za

# Go configuration
GOOS?=$(shell go env GOHOSTOS)
GOARCH?=$(shell go env GOHOSTARCH)

# Add exe extension if windows target
is_windows:=$(filter windows,$(GOOS))
EXT:=$(if $(is_windows),".exe","")

# Archive name
ARCHIVE=$(APPNAME)-$(GOOS)-$(GOARCH).tgz

# Executable name
EXECUTABLE=$(APPNAME)$(EXT)

# Extract version infos
PKG_VERSION:=github.com/ncarlier/$(APPNAME)/pkg/version
VERSION:=`git describe --always --dirty`
GIT_COMMIT:=`git rev-list -1 HEAD --abbrev-commit`
BUILT:=`date`
define LDFLAGS
-X '$(PKG_VERSION).Version=$(VERSION)' \
-X '$(PKG_VERSION).GitCommit=$(GIT_COMMIT)' \
-X '$(PKG_VERSION).Built=$(BUILT)' \
-s -w -buildid=
endef

all: build

# Include common Make tasks
root_dir:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
-include $(root_dir)/.env
makefiles:=$(root_dir)/makefiles
include $(makefiles)/help.Makefile
include $(makefiles)/docker/compose.Makefile

## Clean built files
clean:
	-rm -rf release
	-rm pkg/assets/za.min.js
.PHONY: clean

# Build minified JS
pkg/assets/za.min.js:
	npm run minify

## Build executable
build: pkg/assets/za.min.js
	-mkdir -p release
	echo ">>> Building: $(EXECUTABLE) $(VERSION) for $(GOOS)-$(GOARCH) ..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -tags osusergo,netgo -ldflags "$(LDFLAGS)" -o release/$(EXECUTABLE)
.PHONY: build

release/$(EXECUTABLE): build

## Run tests
test:
	go test ./...
.PHONY: test

## Install executable
install: release/$(EXECUTABLE)
	echo ">>> Installing $(EXECUTABLE) to ${HOME}/.local/bin/$(EXECUTABLE) ..."
	cp release/$(EXECUTABLE) ${HOME}/.local/bin/$(EXECUTABLE)
.PHONY: install

## Create Docker image
image:
	echo ">>> Building Docker image ..."
	docker build --rm -t ncarlier/$(APPNAME) .
.PHONY: image

# Generate changelog
CHANGELOG.md:
	standard-changelog --first-release

var/dbip-country.mmdb:
	echo ">>> Downloading country GeoIP database..."
	mkdir -p var
	wget -O - https://download.db-ip.com/free/dbip-country-lite-$(shell date '+%Y-%m').mmdb.gz | gunzip -c > var/dbip-country-lite.mmdb

## Download Geo IP databases
geoip-db: var/dbip-country.mmdb
.PHONY: geoip-db

## Create archive
archive: release/$(EXECUTABLE)
	echo ">>> Creating release/$(ARCHIVE) archive..."
	tar czf release/$(ARCHIVE) README.md LICENSE CHANGELOG.md -C release/ $(EXECUTABLE)
	rm release/$(EXECUTABLE)
.PHONY: archive

## Create distribution binaries
distribution: CHANGELOG.md
	GOARCH=amd64 make build archive
	GOARCH=arm64 make build archive
	GOARCH=arm make build archive
	GOOS=darwin make build archive
.PHONY: distribution

## Deploy Docker stack
deploy: compose-up 
.PHONY: deploy

## Un-deploy Docker stack
undeploy: compose-down
.PHONY: undeploy
