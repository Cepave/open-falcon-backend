SHELL := /bin/bash
THIS_FILE := $(lastword $(MAKEFILE_LIST))
TARGET_SOURCE = $(shell find main.go g cmd common -name '*.go')
CMD = aggregator graph hbs judge nodata query sender task transfer fe alarm agent mysqlapi f2e-api
TARGET = open-falcon
VERSION := $(shell cat VERSION)

##################################################
# Variables for Different OS information
##################################################
OSTYPE := $(shell echo $$OSTYPE)
# // :~)

##################################################
# For the target of "fmt", "misspell", "fmt-check", and "misspell-check",
# it is possible the listing of files too long to run in shell.
#
# Following variables defines arguments used by "xargs" command
##################################################

# Maximum characters fed to "xargs -s xx"
XARGS_ARGS :=
ifeq ($(findstring cygwin,$(OSTYPE)), cygwin)
	XARGS_ARGS := --max-procs=1 --max-chars=16384
else ifeq ($(findstring linux,$(OSTYPE)), linux)
	XARGS_ARGS := --max-procs=1 --max-chars=65536
else ifeq ($(findstring darwin,$(OSTYPE)), darwin)
	XARGS_ARGS := -P 1 -s 65536
endif

CMD_LIST_GO_FILES := find . -name "*.go" -type f ! -path "./vendor/*"

# Temporary file used to keep listing of go files
LISTFILE_OF_GO_FILES := $(shell mktemp).gofiles.list

XARGS_CMD := xargs $(XARGS_ARGS)

# // :~)

GOFMT ?= gofmt -s

# The folder of GoLang used to search testing package
GO_TEST_FOLDER := common modules scripts/mysql/dbpatch/go
# You should assign the path starting with any of $(GO_TEST_FOLDER)
GO_TEST_EXCLUDE := modules/agent modules/f2e-api modules/fe
# If using verbose
ifeq ($(GO_TEST_VERBOSE), yes)
	run_gotest_verbose = "-v"
endif

all: install $(CMD) $(TARGET)

misspell: build_gofile_listfile .get_misspell
	@echo "Inline fix mis-spelled files.";
	$(XARGS_CMD) misspell -w <$(LISTFILE_OF_GO_FILES);

misspell-check: build_gofile_listfile .get_misspell
	check_cmd="$(XARGS_CMD) misspell -error <$(LISTFILE_OF_GO_FILES)"; \
	echo $$check_cmd; \
	check_output=$$(eval "$$check_cmd"); \
	test -z "$$check_output" || { \
		echo -e "misspell capture error:\n $$check_output\n"; \
		echo "[HELP]" Use \"make misspell\" to fix files inline."(Don't forget to commit changed files)"; \
		exit 1; \
	}

fmt: build_gofile_listfile
	@echo "Inline fix mis-formatted files.";
	$(XARGS_CMD) $(GOFMT) -l -w <$(LISTFILE_OF_GO_FILES);

fmt-check: build_gofile_listfile
	check_cmd="$(XARGS_CMD) $(GOFMT) -d <$(LISTFILE_OF_GO_FILES)"; \
	echo $$check_cmd; \
	check_output=$$(eval "$$check_cmd"); \
	test -z "$$check_output" || { \
		echo -e "gofmt capture error:\n $$check_output\n"; \
		echo "[HELP]" Use \"make fmt\" to fix files inline."(Don't forget to commit changed files)"; \
		exit 1; \
	}

build_gofile_listfile:
	echo Generate "$(LISTFILE_OF_GO_FILES)" file for GoLang files.
	$(CMD_LIST_GO_FILES) >$(LISTFILE_OF_GO_FILES)
	echo -e There are \"`wc -l <$(LISTFILE_OF_GO_FILES)`\" GoLang files."\n"

.get_misspell:
	type -p misspell &>/dev/null || { \
		cmd="go get -v -u github.com/client9/misspell/cmd/misspell"; \
		echo $$cmd; \
		$$cmd; \
	}

go-test:
	./go-test-all.sh -t "$(GO_TEST_FOLDER)" -e "$(GO_TEST_EXCLUDE)" $(run_gotest_verbose)

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
.PHONY: fmt misspell fmt-check misspell-check .get_misspell build_gofile_listfile go-test

.SILENT: build_gofile_listfile misspell-check fmt-check go-test .get_misspell
