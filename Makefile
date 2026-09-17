BINARY := ldap-extractor
CMD := ./cmd
BUILD_DIR := bin
PACKAGE_DIR := bin-package

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
GIT_TAG ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo unknown)
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

VERSION_PKG := ldap-extractor/version

LDFLAGS := \
	-X $(VERSION_PKG).Version=$(VERSION) \
	-X $(VERSION_PKG).GitTag=$(GIT_TAG) \
	-X $(VERSION_PKG).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PKG).BuildTime=$(BUILD_TIME)

.PHONY: all build package run version clean

all: build

build:
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(BINARY)"
	@echo "Version: $(VERSION)"
	@echo "Tag:     $(GIT_TAG)"
	@echo "Commit:  $(GIT_COMMIT)"
	@echo "Built:   $(BUILD_TIME)"
	CGO_ENABLED=0 go build \
		-trimpath \
		-ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY)-$(VERSION) \
		$(CMD)
	
package: build
	@mkdir -p $(PACKAGE_DIR)

	@echo "Packaging $(BINARY)-$(VERSION)"

	rm -rf $(PACKAGE_DIR)/$(BINARY)-$(VERSION)

	mkdir -p $(PACKAGE_DIR)/$(BINARY)-$(VERSION)

	cp $(BUILD_DIR)/$(BINARY)-$(VERSION) \
		$(PACKAGE_DIR)/$(BINARY)-$(VERSION)/

	cp install.sh \
		$(PACKAGE_DIR)/$(BINARY)-$(VERSION)/
	
	cp $(BINARY).service \
		$(PACKAGE_DIR)/$(BINARY)-$(VERSION)/
	
	cp $(BINARY).timer \
		$(PACKAGE_DIR)/$(BINARY)-$(VERSION)/

	tar -C $(PACKAGE_DIR) \
		-czf $(PACKAGE_DIR)/$(BINARY)-$(VERSION).tar.gz \
		$(BINARY)-$(VERSION)

	rm -rf $(PACKAGE_DIR)/$(BINARY)-$(VERSION)

	@echo
	@echo "Package created:"
	@echo "  $(PACKAGE_DIR)/$(BINARY)-$(VERSION).tar.gz"	

run:
	go run $(CMD)

version:
	@echo "Version: $(VERSION)"
	@echo "Tag:     $(GIT_TAG)"
	@echo "Commit:  $(GIT_COMMIT)"
	@echo "Built:   $(BUILD_TIME)"

clean:
	rm -rf $(BUILD_DIR)
	rm -rf $(PACKAGE_DIR)
