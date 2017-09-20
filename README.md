# Open-Falcon Backend

![Open-Falcon](./logo.png)

[![Build Status](https://travis-ci.org/Cepave/open-falcon-backend.svg?branch=develop)](https://travis-ci.org/Cepave/open-falcon-backend)
[![codecov](https://codecov.io/gh/Cepave/open-falcon-backend/branch/develop/graph/badge.svg)](https://codecov.io/gh/Cepave/open-falcon-backend)
[![GoDoc](https://godoc.org/github.com/Cepave/open-falcon-backend?status.svg)](https://godoc.org/github.com/Cepave/open-falcon-backend)
[![Join the chat at https://gitter.im/goappmonitor/Lobby](https://badges.gitter.im/goappmonitor/Lobby.svg)](https://gitter.im/goappmonitor/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Code Health](https://landscape.io/github/Cepave/open-falcon-backend/master/landscape.svg?style=flat)](https://landscape.io/github/Cepave/open-falcon-backend/master)
[![Code Issues](https://www.quantifiedcode.com/api/v1/project/df24b20e9c504ad0a2ac9fa3e99936f5/badge.svg)](https://www.quantifiedcode.com/app/project/df24b20e9c504ad0a2ac9fa3e99936f5)
[![Go Report Card](https://goreportcard.com/badge/github.com/Cepave/open-falcon-backend)](https://goreportcard.com/report/github.com/Cepave/open-falcon-backend)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

# Documentations

- http://book.open-falcon.org
- http://docs.openfalcon.apiary.io

# Get Started

## Start MySQL and Redis using docker

    cd docker
    docker-compose -f init.yml up -d
    docker inspect docker_mysql_1
    docker inspect docker_redis_1
    cd ..

## Change your environment setting

    vi config/confgen.sh

## Start Backend modules

    make clean all pack
    mkdir out
    mv open-falcon-v2.0.0.tar.gz out/
    cd out
    tar zxvf open-falcon-v2.0.0.tar.gz
    ./open-falcon start agent graph transfer hbs fe query

# Compilation

```bash
# all modules
make all

# specified module
make agent
```

# Run Open-Falcon Commands

Agent for example:

    ./open-falcon agent [build|pack|start|stop|restart|status|tail]

# Package Management

We use govendor to manage the golang packages. Please install `govendor` before compilation.

    go get -u github.com/kardianos/govendor

Most depended packages are saved under `./vendor` dir. If you want to add or update a package, just run `govendor fetch xxxx@commitID` or `govendor fetch xxxx@v1.x.x`, then you will find the package have been placed in `./vendor` correctly.

Make sure you're using Go 1.5+ and **GO15VENDOREXPERIMENT=1** env var is exported. (`export GODEBUG=cgocheck=0` using Go 1.6+.)

# Package Release

	make clean all pack

# Testing

## By using `make go-test`

Variables:
* `GO_TEST_FOLDER` - The inclusions of folders(recursively probed) to be tested
* `GO_TEST_EXCLUDE` - The exclusions of folders(include children), which are descendants of `GO_TEST_FOLDER`
* `GO_TEST_VERBOSE` - If the value is "yes", the execution of "go test -test.v"(with additional flags of 3-party frameworks) would be applied

```sh
make go-test GO_TEST_FOLDER="modules common" GO_TEST_EXCLUDE="modules/fe modules/f2e-api"
```

See `Makefile` for default values of the two variables.

## By using `go-test-all.sh`

Arguments:
* `-t` - The inclusions of folders(recursively probed) to be tested
* `-e` - The exclusions of folders(include children), which are descendants of `GO_TEST_FOLDER`
* `-v` - If this flag is shown, the execution of `go test -test.v`(with additional flags of 3-party frameworks) would be applied

```sh
./go-test-all.sh -t "modules common" -e "modules/fe modules/f2e-api"
```
