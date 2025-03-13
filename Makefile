APP=calgo-cli
GOLINT=golint
GODIR=cli
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: cli app

cli: build/bin/$(APP)

build/bin/$(APP): $(GODIR)/*.go
	# mkdir -p build/bin
	cd $(GODIR) && $(GOBUILD) -o ../build/bin/$(APP)

app:
	wails build --platform 'windows'