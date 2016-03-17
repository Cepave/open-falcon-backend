TARGET_SOURCE := $(shell find main.go g commands -name '*.go')
MODULES_SOURCE := $(shell find modules -name '*.go')

CMD = agent aggregator graph hbs judge nodata query sender task transfer fe
BIN = bin/falcon-agent bin/falcon-aggregator bin/falcon-graph bin/falcon-hbs bin/falcon-judge bin/falcon-nodata bin/falcon-query bin/falcon-sender bin/falcon-task bin/falcon-transfer bin/falcon-fe
TARGET = open-falcon

GOTOOLS = github.com/mitchellh/gox golang.org/x/tools/cmd/stringer \
	github.com/jteeuwen/go-bindata/... github.com/elazarl/go-bindata-assetfs/...
PACKAGES=$(shell go list ./... | grep -v '^github.com/Cepave/open-falcon/')
VERSION?=$(shell awk -F\" '/^const Version/ { print $$2; exit }' ./g/version.go)

all: $(BIN) $(TARGET)

$(CMD):
	make bin/falcon-$@

$(TARGET): $(TARGET_SOURCE)
	go build -o open-falcon

$(BIN): $(MODULES_SOURCE)
	go build -o $@ github.com/Cepave/open-falcon/modules/$(@:bin/falcon-%=%)

# dev creates binaries for testing locally - these are put into ./bin and $GOPATH
dev: format
	@CONSUL_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

test: format
	@./scripts/test.sh

format:
	@echo "--> Running go fmt"
	@go fmt $(PACKAGES)

tools:
	go get -u -v $(GOTOOLS)

checkbin:
pack: checkbin
	rm -rf open-falcon-v$(VERSION).tar.gz
	tar -zcvf open-falcon-v$(VERSION).tar.gz ./bin ./config ./open-falcon ./cfg.json

clean:
	git clean -f -d ./bin
	git clean -f -d ./config
	rm -rf ./bin
	rm -rf ./$(TARGET)
	rm -rf open-falcon-v$(VERSION).tar.gz

.PHONY: clean all agent aggregator graph hbs judge nodata query sender task transfer fe
