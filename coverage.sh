#!/usr/bin/env bash

set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor | grep -v \
    -e common/db \
    -e common/gorm \
    -e common/runtime \
    -e common/vipercfg \
    -e scripts/mysql/dbpatch/go/sql \
    \
    -e modules/agent/session \
    -e modules/f2e-api/test \
    -e modules/fe/model/falcon_portal \
    -e modules/fe/test/portal_test \
    -e modules/hbs/http \
    -e modules/query/conf \
    -e modules/query/http \
    ); do
    go get -t $d
    go test -coverprofile=profile.out -covermode=count $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
