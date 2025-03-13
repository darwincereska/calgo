APP=calgo-cli
GOLINT=golint
GODIR=backend
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: app cli

cli:
	cd $(GODIR) && $(GOBUILD) -o ../build/bin/$(APP)

app:
	wails build