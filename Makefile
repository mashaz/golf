GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BINARY_NAME=golf
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_EXE=$(BINARY_NAME).exe


all: clean test build
build:
		$(GOBUILD) -o bin/$(BINARY_NAME) -v
test:
		$(GOTEST) -v 
clean:
		$(GOCLEAN)
		rm -rf bin/*