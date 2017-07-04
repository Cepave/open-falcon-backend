SHELL := /bin/bash
THIS_FILE := $(lastword $(MAKEFILE_LIST))
TARGET_SOURCE = $(shell find main.go g cmd common -name '*.go')
CMD = aggregator graph hbs judge nodata query sender task transfer fe alarm agent nqm-mng f2e-api
TARGET = open-falcon

VERSION := $(shell cat VERSION)

all: install $(CMD) $(TARGET)

$(CMD):
	go get ./modules/$@
	go build -ldflags "-X main.GitCommit=`git log -n1 --pretty=format:%h modules/$@` -X main.Version=${VERSION}" -o bin/$@/falcon-$@ ./modules/$@

$(TARGET): $(TARGET_SOURCE)
	go get .
	go build -ldflags "-X main.GitCommit=`git rev-parse --short HEAD` -X main.Version=$(VERSION)" -o open-falcon

checkvendor:
	@hash govendor > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	@if [ -f ~/.bash_profile ]; then source ~/.bash_profile; fi

install: checkvendor
	govendor sync

checkbin: bin/ config/ open-falcon cfg.json
pack: checkbin
	@if [ -e out ] ; then rm -rf out; fi
	@mkdir out
	@$(foreach var,$(CMD),mkdir -p ./out/$(var)/bin;)
	@$(foreach var,$(CMD),mkdir -p ./out/$(var)/config;)
	@$(foreach var,$(CMD),mkdir -p ./out/$(var)/logs;)
	@$(foreach var,$(CMD),cp ./config/$(var).json ./out/$(var)/config/cfg.json;)
	@$(foreach var,$(CMD),cp ./bin/$(var)/falcon-$(var) ./out/$(var)/bin;)
	@cp -r ./modules/query/js ./modules/query/conf/lambdaSetup.json ./out/query/config
	@cp -r ./modules/fe/{static,views,scripts} ./out/fe/bin
	@cp -r ./modules/alarm/{static,views} ./out/alarm/bin
	@cp -r ./modules/agent/public ./out/agent/bin
	@cp -r ./modules/f2e-api/data ./out/f2e-api/
	@cp cfg.json ./out/cfg.json
	@bash ./config/confgen.sh
	@cp $(TARGET) ./out/$(TARGET)
	tar -C out -zcf open-falcon-v$(VERSION).tar.gz .
	@rm -rf out

coverage: checkvendor
	@$(MAKE) -f $(THIS_FILE) install
	./coverage.sh

clean:
	@rm -rf ./bin
	@rm -rf ./out
	@rm -rf ./$(TARGET)
	@rm -rf ./package_cache_tmp
	@rm -rf open-falcon-v$(VERSION).tar.gz

.PHONY: install clean all aggregator graph hbs judge nodata query sender task transfer fe f2e-api coverage
