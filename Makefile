SHELL := /bin/bash
SOURCE = $(shell find . -name '*.go')
VERSION ?= $(shell awk -F\" '/^const Version/ { print $$2; exit }' ./version.go)

nqm: $(SOURCE)
	go build -o nqm

checkbin: nqm cfg.json
pack: checkbin
	rm -rf nqm-agent-$(VERSION).tar.gz
	tar -zcvf nqm-agent-$(VERSION).tar.gz ./nqm ./cfg.json

clean:
	rm -f ./nqm
	rm -f nqm-agent-$(VERSION).tar.gz
