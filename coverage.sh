#!/usr/bin/env bash

set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor | grep -v \
    -e common/db \
    -e common/gorm \
    -e common/runtime \
    -e common/vipercfg \
    -e modules/agent/session \
    -e scripts/mysql/dbpatch/go/sql \
    -e open-falcon-backend/modules \
    ); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done
