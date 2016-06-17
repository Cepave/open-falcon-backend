SHELL := /bin/bash
TARGET_SOURCE = $(shell find main.go g commands -name '*.go')
CMD = agent aggregator graph hbs judge nodata query sender task transfer fe
BIN = bin/falcon-agent bin/falcon-aggregator bin/falcon-graph bin/falcon-hbs bin/falcon-judge bin/falcon-nodata bin/falcon-query bin/falcon-sender bin/falcon-task bin/falcon-transfer
# bin/falcon-fe
TARGET = open-falcon

GOTOOLS = github.com/mitchellh/gox golang.org/x/tools/cmd/stringer \
	github.com/jteeuwen/go-bindata/... github.com/elazarl/go-bindata-assetfs/...
PACKAGES=$(shell go list ./... | grep -v '^github.com/Cepave/open-falcon/vendor/')
VERSION?=$(shell awk -F\" '/^const Version/ { print $$2; exit }' ./g/version.go)

all: $(BIN) $(TARGET) bin/falcon-fe

$(CMD):
	make bin/falcon-$@

$(TARGET): $(TARGET_SOURCE)
	go build -o open-falcon

$(BIN):
	@cd modules/$(@:bin/falcon-%=%);\
	export commit=`git log -1 --pretty=%h`;\
	echo -e "package g\nconst (\n  COMMIT = \"$$commit\"\n)" > g/git.go
	-@cd vendor/github.com/Cepave/$(@:bin/falcon-%=%);\
	export commit=`git log -1 --pretty=%h`;\
	echo -e "package g\nconst (\n  COMMIT = \"$$commit\"\n)" > g/git.go
	go build -o $@ ./modules/$(@:bin/falcon-%=%)

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

checkbin: bin/ config/ open-falcon cfg.json
pack: checkbin
	rm -rf open-falcon-v$(VERSION).tar.gz
	tar -zcvf open-falcon-v$(VERSION).tar.gz ./bin ./config ./open-falcon ./cfg.json

clean:
	rm -rf ./bin
	rm -rf ./$(TARGET)
	rm -rf open-falcon-v$(VERSION).tar.gz

.PHONY: clean all agent aggregator graph hbs judge nodata query sender task transfer fe

bin/falcon-agent : $(shell find modules/agent/ -name '*.go')
bin/falcon-aggregator : $(shell find modules/aggregator/ -name '*.go')
bin/falcon-graph : $(shell find modules/graph/ -name '*.go')
bin/falcon-hbs : $(shell find modules/hbs/ -name '*.go')
bin/falcon-judge : $(shell find modules/judge/ -name '*.go')
bin/falcon-nodata : $(shell find modules/nodata/ -name '*.go')
bin/falcon-query : $(shell find modules/query/ -name '*.go')
	go build -o $@ ./modules/query
	cp -r modules/query/js bin/js
	mkdir -p bin/conf
	cp modules/query/conf/lambdaSetup.json bin/conf
bin/falcon-sender : $(shell find modules/sender/ -name '*.go')
bin/falcon-task : $(shell find modules/task/ -name '*.go')
bin/falcon-transfer : $(shell find modules/transfer/ -name '*.go')
bin/falcon-fe: $(shell find modules/fe/ -name '*.go')
	go build -o $@ ./modules/$(@:bin/falcon-%=%)
	mkdir -p bin/fe
	cp -r modules/fe/{control,cfg.example.json,conf,static,views,scripts} bin/fe/
	cp bin/falcon-fe bin/fe/falcon-fe
