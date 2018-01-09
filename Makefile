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
GO_TEST_FOLDER := common modules
# You should assign the path starting with any of $(GO_TEST_FOLDER)
GO_TEST_EXCLUDE := modules/agent modules/f2e-api modules/fe

GO_TEST_FLAGS :=
GO_TEST_VERBOSE :=
GO_TEST_PROPS :=
mp =
GO_TEST_PROPS_SEP :=
GO_TEST_PROPS_FILE :=
GO_TEST_COVERAGE_FILE :=

all: setup-govendor govendor-check $(CMD) $(TARGET)

check-all: misspell-check fmt-check govendor-check

misspell: build_gofile_listfile .get_misspell
	@echo "Inline fix mis-spelled files.";
	$(XARGS_CMD) misspell -w <$(LISTFILE_OF_GO_FILES);

misspell-check: build_gofile_listfile .get_misspell
	check_cmd="$(XARGS_CMD) misspell -error <$(LISTFILE_OF_GO_FILES)"; \
	echo Check misspelling of GoLang ...; \
	echo -e "\t$$check_cmd"; \
	check_output=$$(eval "$$check_cmd"); \
	test -z "$$check_output" || { \
		echo -e "misspell capture error:\n $$check_output\n"; \
		echo "[HELP]" Use \"make misspell\" to fix files inline."(Don't forget to commit changed files)"; \
		exit 1; \
	}
	echo -e "[PASS]\n"

fmt: build_gofile_listfile
	@echo "Inline fix mis-formatted files.";
	$(XARGS_CMD) $(GOFMT) -l -w <$(LISTFILE_OF_GO_FILES);

fmt-check: build_gofile_listfile
	check_cmd="$(XARGS_CMD) $(GOFMT) -d <$(LISTFILE_OF_GO_FILES)"; \
	echo Check formatter of GoLang ...; \
	echo -e "\t$$check_cmd"; \
	check_output=$$(eval "$$check_cmd"); \
	test -z "$$check_output" || { \
		echo -e "gofmt capture error:\n $$check_output\n"; \
		echo "[HELP]" Use \"make fmt\" to fix files inline."(Don't forget to commit changed files)"; \
		exit 1; \
	};
	echo -e "[PASS]\n"

# Asserts that there is no external libaray
govendor-check:
	@echo -n "Check +external of govendor(\"govendor list +external\") ... "
	@external_libs=`govendor list +external`; \
	if test -n "$$external_libs"; then \
		echo -e "There are external library. You should use govendor add \"<lib path>\" .\n"; \
		echo -e $$external_libs; \
		exit 1; \
	fi
	@echo "[PASS]"
	@echo -n "Check +unused of govendor(\"govendor list +unused\") ... "
	@unused_libs=`govendor list +unused`; \
	if test -n "$$unused_libs"; then \
		echo -e "There are unused library. You should use govendor remove \"<lib path>\" .\n"; \
		echo -e $$unused_libs; \
		exit 1; \
	fi
	@echo "[PASS]"

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
	./go-test-all.sh \
		-t "$(GO_TEST_FOLDER)" -e "$(GO_TEST_EXCLUDE)" \
		$(if $(strip $(GO_TEST_COVERAGE_FILE)),-c "$(GO_TEST_COVERAGE_FILE)",) \
		$(if $(strip $(GO_TEST_PROPS_FILE)),-f "$(GO_TEST_PROPS_FILE)",) \
		$(if $(strip $(GO_TEST_PROPS)),-p "$(GO_TEST_PROPS)",) \
		$(if $(strip $(GO_TEST_PROPS_SEP)),-s "$(GO_TEST_PROPS_SEP)",) \
		$(if $(strip $(GO_TEST_FLAGS)),-a "$(GO_TEST_FLAGS)",) \
		$(if $(filter yes,$(GO_TEST_VERBOSE)),-v,)

$(CMD):
	go build -ldflags "-X main.GitCommit=`git log -n1 --pretty=format:%h modules/$@` -X main.Version=${VERSION}" -o bin/$@/falcon-$@ ./modules/$@

$(TARGET): $(TARGET_SOURCE)
	go build -ldflags "-X main.GitCommit=`git rev-parse --short HEAD` -X main.Version=$(VERSION)" -o $@

setup-govendor:
	@hash govendor > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/kardianos/govendor; \
	fi
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

clean:
	@rm -rf ./bin
	@rm -rf ./out
	@rm -rf ./$(TARGET)
	@rm -rf open-falcon-v$(VERSION).tar.g

.PHONY: install clean all aggregator graph hbs judge nodata query sender task transfer fe f2e-api coverage
.PHONY: fmt misspell fmt-check misspell-check check-all .get_misspell build_gofile_listfile go-test

.SILENT: build_gofile_listfile misspell-check fmt-check go-test .get_misspell
