BUILDPATH=$(CURDIR)
GO=$(shell which go)
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
EXEBUILD=main.go
EXEFILE=main
EXENAME=$(notdir $(CURDIR))
USER1=$(shell users | cut -f 1 -d ' ')
OS := $(shell uname)
ARCH := $(shell uname -m)

ifeq ($(OS), Linux)
    GOOS := linux
    GOARCH := amd64
    EXEEXT :=
    INSTALLPATH := /usr/local/bin
endif

ifeq ($(OS), Darwin)
    GOOS := darwin
    ifeq ($(ARCH), x86_64)
        GOARCH := amd64
    else ifeq ($(ARCH), arm64)
        GOARCH := arm64
    endif
    EXEEXT :=
    INSTALLPATH := /usr/local/bin
endif

ifeq ($(OS), CYGWIN)
    GOOS := windows
    GOARCH := amd64
    EXEEXT := .exe
    INSTALLPATH := C:/Program\ Files/YourApp
endif

ifeq ($(OS), MINGW)
    GOOS := windows
    GOARCH := amd64
    EXEEXT := .exe
    INSTALLPATH := C:/Program\ Files/YourApp
endif

ifeq ($(OS), MSYS)
    GOOS := windows
    GOARCH := amd64
    EXEEXT := .exe
    INSTALLPATH := C:/Program\ Files/YourApp
endif

echo:
	@echo "Build Path: $(BUILDPATH)"
	@echo "Go Path: $(GO)"
	@echo "Go Build: $(GOBUILD)"
	@echo "Go Clean: $(GOCLEAN)"
	@echo "Executable Build File: $(EXEBUILD)"
	@echo "Executable File: $(EXEFILE)"
	@echo "New Executable Name: $(EXENAME)"
	@echo "Executable Extension: $(EXEEXT)"
	@echo "Target OS: $(GOOS)"
	@echo "Target Arch: $(GOARCH)"
	@echo "Install Path: $(INSTALLPATH)"
	@echo "User: $(USER1)"

makedir:
	export GOPATH=$(CURDIR)
	@echo "Start building tree..."
	@if [ ! -d $(BUILDPATH)/bin ] ; then mkdir -pv $(BUILDPATH)/bin; fi
	@if [ ! -d $(BUILDPATH)/pkg ] ; then mkdir -pv $(BUILDPATH)/pkg; fi

build: makedir
	export GOPATH=$(CURDIR)
	@echo "Start building executable for $(GOOS)/$(GOARCH)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 $(GOBUILD) -o ./bin/$(EXENAME)-$(GOOS)-$(GOARCH)$(EXEEXT)
	@echo "Build completed for $(GOOS)/$(GOARCH)..."

install: build
	@echo "Installing executable..."
	@if [ "$(GOOS)" = "windows" ]; then \
		cp -v ./bin/$(EXENAME)-$(GOOS)-$(GOARCH)$(EXEEXT) $(INSTALLPATH)/$(EXENAME)-$(GOOS)-$(GOARCH)$(EXEEXT); else \
		sudo cp -v ./bin/$(EXENAME)-$(GOOS)-$(GOARCH)$(EXEEXT) $(INSTALLPATH)/$(EXENAME) \
		&& sudo chown -v root:$(USER1) $(INSTALLPATH)/$(EXENAME) \
		&& sudo chmod -v 550 $(INSTALLPATH)/$(EXENAME); fi
	@echo "Install completed."
	@make clean

clean:
	@echo "Cleaning up..."
	@rm -vrf $(BUILDPATH)/bin/*
	@rm -vrf $(BUILDPATH)/pkg/*

# 	@rm -vrf $(BUILDPATH)/bin/$(EXENAME)-$(GOOS)-$(GOARCH)$(EXEEXT)
# 	@rm -vrf $(BUILDPATH)/pkg

# // GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o dragonxf_intel
# // GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dragonxf_linux
# // GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o dragonxf_arm
