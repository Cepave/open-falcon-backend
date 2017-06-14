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
	@echo ">> Sync vendor files before testing code coverage"
	govendor fetch github.com/jpillora/backoff@06c7a16c845dc8e0bf575fafeeca0f5462f5eb4d
	govendor fetch github.com/Pallinder/go-randomdata@8c3362a5e6781e0d1046f1267b3c1f19b2cde334
	govendor fetch github.com/chyeh/pubip@b7e679cf541cd580e99712413b9351008731d09d
	govendor fetch github.com/montanaflynn/stats@41c34e4914ec3c05d485e564d9028d8861d5d9ad
	govendor fetch github.com/parnurzeal/gorequest@5bf13be198787abbed057fb7c4007f372083a0f5
	govendor fetch github.com/patrickmn/go-cache@7ac151875ffb48b9f3ccce9ea20f020b0c1596c8
	govendor fetch github.com/smartystreets/goconvey/convey@af8e7d560364b90f732a1d119d17b5506e50447d
	govendor fetch gopkg.in/check.v1@20d25e2804050c1cd24a7eea1e7a6447dd0e74ec
	govendor fetch github.com/onsi/ginkgo@502bce873ec80059e9465bf32a9aa61e891b7009
	govendor fetch github.com/onsi/ginkgo/extensions/table@502bce873ec80059e9465bf32a9aa61e891b7009
	govendor fetch github.com/onsi/gomega@00acfa9d92a386415bd235ab069c52063f925998
	@$(MAKE) -f $(THIS_FILE) install
	./coverage.sh

clean:
	@rm -rf ./bin
	@rm -rf ./out
	@rm -rf ./$(TARGET)
	@rm -rf ./package_cache_tmp
	@rm -rf open-falcon-v$(VERSION).tar.gz

.PHONY: install clean all aggregator graph hbs judge nodata query sender task transfer fe f2e-api coverage
