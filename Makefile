SHELL := /bin/bash
TARGET_SOURCE = $(shell find main.go g commands -name '*.go')
CMD = aggregator graph hbs judge nodata query sender task transfer fe
BIN = bin/falcon-aggregator bin/falcon-graph bin/falcon-hbs bin/falcon-judge bin/falcon-nodata bin/falcon-sender bin/falcon-task bin/falcon-transfer
TARGET = open-falcon

PACKAGES=$(shell go list ./... | grep -v '^github.com/Cepave/open-falcon/vendor/')
VERSION?=$(shell awk -F\" '/^const Version/ { print $$2; exit }' ./g/version.go)

all: $(BIN) $(TARGET) bin/falcon-fe bin/falcon-query

$(CMD):
	make bin/falcon-$@

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
	@go fmt $(PACKAGES)

checkbin: bin/ config/ open-falcon cfg.json
pack: checkbin
	mkdir out
	$(foreach var,$(CMD),mkdir -p ./out/$(var)/bin;)
	$(foreach var,$(CMD),mkdir -p ./out/$(var)/config;)
	$(foreach var,$(CMD),mkdir -p ./out/$(var)/logs;)
	$(foreach var,$(CMD),cp ./config/$(var).json ./out/$(var)/config/cfg.json;)
	$(foreach var,$(CMD),cp ./bin/falcon-$(var) ./out/$(var)/bin;)
	cp cfg.json ./out/cfg.json
	cp $(TARGET) ./out/$(TARGET)
	cd out && tar -zcvf open-falcon-v$(VERSION).tar.gz ./*
	cd ..
	mv ./out/open-falcon-v$(VERSION).tar.gz .
	rm -rf out

clean:
	rm -rf ./bin
	rm -rf ./out
	rm -rf ./$(TARGET)
	rm -rf open-falcon-v$(VERSION).tar.gz

.PHONY: clean all aggregator graph hbs judge nodata query sender task transfer fe

bin/falcon-agent : $(shell find modules/agent/ -name '*.go')
bin/falcon-nqm-agent : $(shell find modules/nqm-agent/ -name '*.go')
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
