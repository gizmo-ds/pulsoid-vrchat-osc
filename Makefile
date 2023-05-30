NAME=pulsoid-vrchat-osc
MAIN=cmd/cli/main.go
PKGNAME=github.com/gizmo-ds/pulsoid-vrchat-osc
OUTDIR=build
VERSION=$(shell git describe --tags --always --dirty)
FLAGS+=-trimpath
FLAGS+=-tags timetzdata
FLAGS+=-ldflags "-s -w -X main.AppVersion=$(VERSION)"
export CGO_ENABLED=0

all: windows-amd64

initialize:
	mkdir -p build
	cp config_example.toml build/config.toml
	cp README*.md LICENSE build

generate:
	go generate ./...

windows-amd64: initialize generate
	GOOS=windows GOARCH=amd64 go build $(FLAGS) -o $(OUTDIR)/$(NAME).exe $(MAIN)
	cd $(OUTDIR) && zip $(NAME)-$@.zip ./*

sha256sum:
	cd $(OUTDIR); for file in *.zip; do sha256sum $$file > $$file.sha256; done

clean:
	rm -rf $(OUTDIR)/*
