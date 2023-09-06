PROJ_NAME = gndiff

VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S%Z')

FLAGS_INTEL64 = GOARCH=amd64
NO_C = CGO_ENABLED=0
FLAGS_LINUX = $(FLAGS_INTEL64) GOOS=linux
FLAGS_MAC = $(FLAGS_INTEL64) GOOS=darwin
FLAGS_MAC_ARM = $GOARCH=arm64 GOOS=darwin
FLAGS_WIN = $(FLAGS_INTEL64) GOOS=windows
FLAGS_LD = -ldflags "-X github.com/gnames/$(PROJ_NAME)/pkg.Build=$(DATE) \
                  -X github.com/gnames/$(PROJ_NAME)/pkg.Version=$(VERSION)"
FLAGS_REL = -trimpath -ldflags "-s -w \
						-X github.com/gnames/$(PROJ_NAME)/pkg.Build=$(DATE)"
GOCMD = go
GOBUILD = $(GOCMD) build $(FLAGS_LD)
GORELEASE = $(GOCMD) build $(FLAGS_REL)
GOINSTALL = $(GOCMD) install $(FLAGS_LD)
GOCLEAN = $(GOCMD) clean
GOGET = $(GOCMD) get

RELEASE_DIR ?= "/tmp"
BUILD_DIR ?= "."

all: install

test: deps install
	@echo Run tests
	$(GOCMD) test -shuffle=on -count=1 -race -coverprofile=coverage.txt ./...

test-build: deps build

deps:
	@echo Download go.mod dependencies
	$(GOCMD) mod download;

tools: deps
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

build:
	@echo Building
	$(GOCLEAN); \
	$(FLAGS_SHARED) $(NO_C) $(GOBUILD) -o $(BUILD_DIR)

buildrel:
	@echo Building release binary
	$(GOCLEAN); \
	$(NO_C) $(GORELEASE) -o $(BUILD_DIR);

install:
	$(GOCLEAN); \
	$(FLAGS_SHARED) $(NO_C) $(GOINSTALL)

release: dockerhub
	@echo Make release
	$(GOCLEAN); \
	$(FLAGS_LINUX) $(NO_C) $(GOBUILD); \
	tar zcf $(RELEASE_DIR)/$(PROJ_NAME)-$(VER)-linux.tar.gz $(PROJ_NAME); \
	$(GOCLEAN); \
	$(FLAGS_MAC) $(NO_C) $(GOBUILD); \
	tar zcf $(RELEASE_DIR)/$(PROJ_NAME)-$(VER)-mac.tar.gz $(PROJ_NAME); \
	$(GOCLEAN); \
	$(FLAGS_MAC_ARM) $(NO_C) $(GOBUILD); \
	tar zcf $(RELEASE_DIR)/$(PROJ_NAME)-$(VER)-mac-arm64.tar.gz $(PROJ_NAME); \
	$(GOCLEAN); \
	$(FLAGS_WIN) $(NO_C) $(GOBUILD); \
	zip -9 $(RELEASE_DIR)/$(PROJ_NAME)-$(VER)-win-64.zip $(PROJ_NAME).exe; \
	$(GOCLEAN);

docker: buildrel
	@echo Build Docker images
	docker buildx build -t gnames/$(PROJ_NAME):latest -t gnames/$(PROJ_NAME):$(VERSION) .; \
	$(GOCLEAN);

dockerhub: docker
	@echo Push Docker images to DockerHub
	docker push gnames/$(PROJ_NAME); \
	docker push gnames/$(PROJ_NAME):$(VERSION)
