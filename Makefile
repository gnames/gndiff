VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S%Z')

FLAGS_INTEL64 = GOARCH=amd64
NO_C = CGO_ENABLED=0
FLAGS_LINUX = $(FLAGS_INTEL) GOOS=linux
FLAGS_MAC = $(FLAGS_INTEL) GOOS=darwin
FLAGS_MAC_ARM = $GOARCH=arm64 GOOS=darwin
FLAGS_WIN = $(FLAGS_INTEL64) GOOS=windows
FLAGS_LD=-ldflags "-s -w -X github.com/gnames/gndiff.Build=${DATE} \
                  -X github.com/gnames/gndiff.Version=${VERSION}"
GOCMD = go
GOBUILD = $(GOCMD) build $(FLAGS_LD)
GOINSTALL = $(GOCMD) install $(FLAGS_LD)
GOCLEAN = $(GOCMD) clean
GOGET = $(GOCMD) get

RELEASE_DIR ?= "/tmp"
BUILD_DIR ?= "."
CLIB_DIR ?= "."

all: install

test: deps install
	$(FLAG_INTEL) go test -race ./...

test-build: deps build

deps:
	$(GOCMD) mod download;

tools: deps
	@echo Installing tools from tools.go
	@cat gndiff/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

build:
	cd gndiff; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) $(NO_C) $(GOBUILD) -o $(BUILD_DIR)

install:
	cd gndiff; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) $(NO_C) $(GOINSTALL)

release: dockerhub
	cd gndiff; \
	$(GOCLEAN); \
	$(FLAGS_LINUX) $(NO_C) $(GOBUILD); \
	tar zcf $(RELEASE_DIR)/gndiff-$(VER)-linux.tar.gz gndiff; \
	$(GOCLEAN); \
	$(FLAGS_MAC) $(NO_C) $(GOBUILD); \
	tar zcf $(RELEASE_DIR)/gndiff-$(VER)-mac.tar.gz gndiff; \
	$(GOCLEAN); \
	$(FLAGS_MAC_ARM) $(NO_C) $(GOBUILD); \
	tar zcf $(RELEASE_DIR)/gndiff-$(VER)-mac-arm64.tar.gz gndiff; \
	$(GOCLEAN); \
	$(FLAGS_WIN) $(NO_C) $(GOBUILD); \
	zip -9 $(RELEASE_DIR)/gndiff-$(VER)-win-64.zip gndiff.exe; \
	$(GOCLEAN);

dc: build
	docker-compose build;

docker: build
	docker build -t gnames/gndiff:latest -t gnames/gndiff:$(VERSION) .; \
	cd gndiff; \
	$(GOCLEAN);

dockerhub: docker
	docker push gnames/gndiff; \
	docker push gnames/gndiff:$(VERSION)
