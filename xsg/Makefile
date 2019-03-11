GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean 
BINARY_NAME=build/xsg

all: clean go-install go-build 
go-install: 
	$(GOCMD) get
go-build:
	$(GOBUILD) -o $(BINARY_NAME)
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

