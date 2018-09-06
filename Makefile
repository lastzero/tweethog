export GO111MODULE=on
GOCMD=go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
BINARY_NAME=tweethog

all: deps build
install:
	$(GOINSTALL) cmd/$(BINARY_NAME)/$(BINARY_NAME).go
build:
	$(GOBUILD) cmd/$(BINARY_NAME)/$(BINARY_NAME).go
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
format:
	$(GOFMT) ./...
deps:
	$(GOBUILD) -v ./...
upgrade:
	$(GOGET) -u