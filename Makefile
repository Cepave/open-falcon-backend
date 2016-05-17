SHELL := /bin/bash
SOURCE = $(shell find . -name '*.go')
VERSION ?= $(shell awk -F\" '/^var Version/ { print $$2; exit }' ./version.go)

nqm-agent: $(SOURCE)
	go build -o nqm-agent

checkbin: nqm-agent cfg.json
pack: checkbin
	rm -rf nqm-agent-$(VERSION).tar.gz
	tar -zcvf nqm-agent-$(VERSION).tar.gz ./nqm-agent ./cfg.json

clean:
	rm -f ./nqm-agent
	rm -f nqm-agent-$(VERSION).tar.gz
