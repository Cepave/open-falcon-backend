SHELL := /bin/bash
TARGET_SOURCE = $(shell find main.go g commands -name '*.go')
CMD = aggregator graph hbs judge nodata query sender task transfer fe
BIN = bin/falcon-aggregator bin/falcon-graph bin/falcon-hbs bin/falcon-judge bin/falcon-nodata bin/falcon-sender bin/falcon-task bin/falcon-transfer
TARGET = open-falcon

VERSION?=$(shell awk -F\" '/^const Version/ { print $$2; exit }' ./g/version.go)

all: $(BIN) $(TARGET) bin/falcon-fe bin/falcon-query

$(CMD):
	@make bin/falcon-$@

$(TARGET): $(TARGET_SOURCE)
	go build -o open-falcon

$(BIN):
	go build -o $@ ./modules/$(@:bin/falcon-%=%)

# dev creates binaries for testing locally - these are put into ./bin and $GOPATH
dev: format
	@CONSUL_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

test: format
	@./scripts/test.sh

format:
	@echo "--> Running go fmt"
	@go fmt `go list ./...`

checkbin: bin/ config/ open-falcon cfg.json
pack: checkbin
	@if [ -e out ] ; then rm -rf out; fi
	@mkdir out
	@$(foreach var,$(CMD),mkdir -p ./out/$(var)/bin;)
	@$(foreach var,$(CMD),mkdir -p ./out/$(var)/config;)
	@$(foreach var,$(CMD),mkdir -p ./out/$(var)/logs;)
	@$(foreach var,$(CMD),cp ./config/$(var).json ./out/$(var)/config/cfg.json;)
	@$(foreach var,$(CMD),cp ./bin/falcon-$(var) ./out/$(var)/bin;)
	@cp cfg.json ./out/cfg.json
	@cp $(TARGET) ./out/$(TARGET)
	tar -C out -zcf open-falcon-v$(VERSION).tar.gz .
	@rm -rf out

clean:
	@rm -rf ./bin
	@rm -rf ./out
	@rm -rf ./$(TARGET)
	@rm -rf open-falcon-v$(VERSION).tar.gz

.PHONY: clean all aggregator graph hbs judge nodata query sender task transfer fe

bin/falcon-query:
	go build -o $@ ./modules/query
	cp -r modules/query/js bin/js
	mkdir -p bin/conf
	cp modules/query/conf/lambdaSetup.json bin/conf

bin/falcon-fe:
	go build -o $@ ./modules/$(@:bin/falcon-%=%)
	mkdir -p bin/fe
	cp -r modules/fe/{control,cfg.example.json,conf,static,views,scripts} bin/fe/
	cp bin/falcon-fe bin/fe/falcon-fe
