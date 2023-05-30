NAME=example
MAIN=cmd/example/main.go
PKGNAME=github.com/gizmo-ds/go-cli-template
OUTDIR=build
VERSION=$(shell git describe --tags --always --dirty)
FLAGS+=-trimpath
FLAGS+=-tags timetzdata
FLAGS+=-ldflags "-s -w -X main.AppVersion=$(VERSION)"
export CGO_ENABLED=0

all: windows-amd64 linux-amd64 darwin-amd64

generate:
	go generate ./...

darwin-amd64: generate
	GOOS=darwin GOARCH=amd64 go build $(FLAGS) -o $(OUTDIR)/$(NAME)-$@ $(MAIN)

linux-amd64: generate
	GOOS=linux GOARCH=amd64 go build $(FLAGS) -o $(OUTDIR)/$(NAME)-$@ $(MAIN)

windows-amd64: generate
	GOOS=windows GOARCH=amd64 go build $(FLAGS) -o $(OUTDIR)/$(NAME)-$@.exe $(MAIN)

sha256sum:
	cd $(OUTDIR); for file in *; do sha256sum $$file > $$file.sha256; done