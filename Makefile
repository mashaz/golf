GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BINARY_NAME=golf
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_EXE=$(BINARY_NAME).exe


all: clean test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v
test:
		$(GOTEST) -v 
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
		rm -f $(BINARY_EXE)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v -h
		./$(BINARY_NAME)