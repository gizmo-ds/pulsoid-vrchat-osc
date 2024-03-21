NAME=pulsoid-vrchat-osc
MAIN=cmd/cli/main.go
PKGNAME=github.com/gizmo-ds/pulsoid-vrchat-osc
OUTDIR=build
VERSION=$(shell git describe --tags --always --dirty)
FLAGS+=-trimpath
FLAGS+=-tags timetzdata
FLAGS+=-ldflags "-s -w -X $(PKGNAME)/internal/global.AppVersion=$(VERSION)"
export CGO_ENABLED=0

PLATFORMS := windows

all: build-all compress zip sha256sum

initialize:
	@mkdir -p $(OUTDIR)
	@cp config_example.toml $(OUTDIR)/config.toml
	@cp README*.md LICENSE $(OUTDIR)

build-all: $(PLATFORMS)

generate:
	go generate ./...

$(PLATFORMS): generate
	GOOS=$@ GOARCH=amd64 go build $(FLAGS) -o $(OUTDIR)/$(NAME)-$@$(if $(filter windows,$@),.exe) $(MAIN)

sha256sum: zip
	@cd $(OUTDIR); sha256sum *.zip > sha256.txt

zip:
	@cp config_example.toml $(OUTDIR)/config.toml
	for platform in $(PLATFORMS); do \
		zip -jq9 $(OUTDIR)/$(NAME)-$$platform.zip $(OUTDIR)/$(NAME)-$$platform* $(OUTDIR)/config.toml README*.md LICENSE; \
	done

compress:
	@if [ -n "$(shell command -v upx 2> /dev/null)" ]; then for file in build/*.exe; do upx $$file; done; fi

clean:
	@rm -rf $(OUTDIR)/*
