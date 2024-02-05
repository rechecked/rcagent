GOCMD=go
GOTEST=$(GOCMD) test
VERSION?=$(shell cat VERSION | tr -d "\n")
RELEASE?=1
BINARY_NAME=rcagent
DIR_NAME=rcagent-$(VERSION)
LOCAL_DIR=/usr/local/rcagent
VFLAGS=-X github.com/rechecked/rcagent/internal/config.Version=$(VERSION)
LDFLAGS?=

.PHONY: build test clean install

all: help

build: clean
	$(GOCMD) build -o build/bin/$(BINARY_NAME) -ldflags "$(VFLAGS) $(LDFLAGS)"

build-tar: clean
	tar -czf build/rcagent-$(VERSION).tar.gz . --transform 's,^,rcagent-$(VERSION)/,'

build-rpm: build-tar
	mkdir -p $(HOME)/rpmbuild/SOURCES/
	mv -f build/rcagent-$(VERSION).tar.gz $(HOME)/rpmbuild/SOURCES/
	cp build/package/rcagent.spec build/rcagent.spec
	sed -i "s/Version:.*/Version: $(VERSION)/g" build/rcagent.spec
	sed -i "s/Release:.*/Release: $(RELEASE)%{?dist}/g" build/rcagent.spec
	rpmbuild -ba build/rcagent.spec
	find $(HOME)/rpmbuild/RPMS -name "rcagent-*.rpm" -exec cp {} build \;

build-deb: build-rpm
	cd build && alien -c -k -v rcagent-*.rpm

build-dmg:
	mkdir build/$(DIR_NAME)
	cp build/bin/$(BINARY_NAME) build/$(DIR_NAME)/$(BINARY_NAME)
	cp build/package/config.yml build/$(DIR_NAME)/config.yml
	cp build/package/macos/install.sh build/$(DIR_NAME)/install.sh
	cp build/package/macos/uninstall.sh build/$(DIR_NAME)/uninstall.sh
	cd build && hdiutil create -volname $(DIR_NAME) -srcfolder $(DIR_NAME) -ov -format UDZO $(DIR_NAME).dmg

install:
	mkdir -p $(LOCAL_DIR)/plugins
	cp -f build/bin/$(BINARY_NAME) $(LOCAL_DIR)/$(BINARY_NAME)
	cp -n build/package/config.yml $(LOCAL_DIR)/config.yml

test:
	$(GOTEST) -v ./... -coverprofile cover.out

coverage:
	$(GOCMD) tool cover -func cover.out

clean:
	rm -rf build/bin
	rm -rf build/rcagent-*
	rm -f build/rcagent.spec

help:
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@echo '  build 			build the binary'
	@echo ''
	@echo '  install        install rcagent into /usr/local'
	@echo ''
	@echo '  build-rpm		build rpm package'
	@echo '  build-deb		build deb package'
	@echo '  build-dmg		build dmg package'
	@echo ''
	@echo '  build-tar		bundle the source into a tar.gz file'
	@echo ''
	@echo '  test 			run the go tests'
	@echo '  coverage       show the coverage from running make test'
	@echo '  clean			clean up the directoies/binary files'
