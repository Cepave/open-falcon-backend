SHELL := /bin/bash
THIS_FILE := $(lastword $(MAKEFILE_LIST))
TARGET_SOURCE = $(shell find main.go g cmd common -name '*.go')
CMD = aggregator graph hbs judge nodata query sender task transfer fe alarm agent mysqlapi f2e-api
TARGET = open-falcon
VERSION := $(shell cat VERSION)

##################################################
# For the target of "fmt", "misspell", "fmt-check", and "misspell-check",
# it is possible the listing of files too long to run in shell.
#
# Following variables defines arguments used by "xargs" command
##################################################

# Maximum characters fed to "xargs -s xx"
CMD_MAX_CHARS := 16384
CMD_LIST_GO_FILES := find . -name "*.go" -type f ! -path "./vendor/*"

# Temporary file used to keep listing of go files
LISTFILE_OF_GO_FILES := $(shell mktemp).gofiles.list

XARGS_CMD := xargs --max-procs=1 -s $(CMD_MAX_CHARS) --arg-file=$(LISTFILE_OF_GO_FILES)

# // :~)

GOFMT ?= gofmt -s

# The folder of GoLang used to search testing package
GO_TEST_FOLDER := common modules scripts/mysql/dbpatch/go
# You should assign the path starting with any of $(GO_TEST_FOLDER)
GO_TEST_EXCLUDE := modules/agent modules/f2e-api modules/fe

all: install $(CMD) $(TARGET)

fmt: build_gofile_listfile
	$(XARGS_CMD) $(GOFMT) -w

misspell: build_gofile_listfile
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	$(XARGS_CMD) misspell -w

fmt-check: build_gofile_listfile
	# get all go files and run go fmt on them
	@diff=$$($(XARGS_CMD) $(GOFMT) -d); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

misspell-check: build_gofile_listfile
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	$(XARGS_CMD) misspell -error

build_gofile_listfile:
	echo Generate "$(LISTFILE_OF_GO_FILES)" file for GoLang files.
	$(CMD_LIST_GO_FILES) >$(LISTFILE_OF_GO_FILES)
	echo There are \"`wc -l <$(LISTFILE_OF_GO_FILES)`\" GoLang files.

go-test:
	./go-test-all.sh -t "$(GO_TEST_FOLDER)" -e "$(GO_TEST_EXCLUDE)"

$(CMD):
	go build -ldflags "-X main.GitCommit=`git log -n1 --pretty=format:%h modules/$@` -X main.Version=${VERSION}" -o bin/$@/falcon-$@ ./modules/$@

$(TARGET): $(TARGET_SOURCE)
	go build -ldflags "-X main.GitCommit=`git rev-parse --short HEAD` -X main.Version=$(VERSION)" -o $@

checkvendor:
	@hash govendor > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/kardianos/govendor; \
	fi

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
	@cp -r ./modules/f2e-api/data ./out/f2e-api/bin
	@cp -r ./modules/f2e-api/lambda_extends/js ./out/f2e-api/bin
	@cp -r ./modules/f2e-api/lambda_extends/conf/lambdaSetup.json ./out/f2e-api/config
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
	@rm -rf open-falcon-v$(VERSION).tar.g

.PHONY: install clean all aggregator graph hbs judge nodata query sender task transfer fe f2e-api coverage
.PHONY: fmt misspell fmt-check misspell-check build_gofile_listfile go-test

.SILENT: build_gofile_listfile go-test
