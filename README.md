# Open-Falcon Backend

![Open-Falcon](./logo.png)

[![Build Status](https://travis-ci.org/Cepave/open-falcon-backend.svg?branch=master)](https://travis-ci.org/Cepave/open-falcon-backend)
[![codecov](https://codecov.io/gh/Cepave/open-falcon-backend/branch/master/graph/badge.svg)](https://codecov.io/gh/Cepave/open-falcon-backend)
[![GoDoc](https://godoc.org/github.com/Cepave/open-falcon-backend?status.svg)](https://godoc.org/github.com/Cepave/open-falcon-backend)
[![Join the chat at https://gitter.im/goappmonitor/Lobby](https://badges.gitter.im/goappmonitor/Lobby.svg)](https://gitter.im/goappmonitor/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Code Health](https://landscape.io/github/Cepave/open-falcon-backend/master/landscape.svg?style=flat)](https://landscape.io/github/Cepave/open-falcon-backend/master)
[![Code Issues](https://www.quantifiedcode.com/api/v1/project/98b2cb0efd774c5fa8f9299c4f96a8c5/badge.svg)](https://www.quantifiedcode.com/app/project/98b2cb0efd774c5fa8f9299c4f96a8c5)
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
## How-to

Make sure you're using Go 1.5+ and **GO15VENDOREXPERIMENT=1** env var is exported. (`export GODEBUG=cgocheck=0` using Go 1.6+.)

 0. Install `trash` by `go get github.com/rancher/trash`.
 1. Edit `trash.yml` file to your needs. See the example as follow.
 2. Run `trash --keep` to download the dependencies.

```yaml
package: github.com/Cepave/open-falcon-backend

import:
- package: github.com/Cepave/common              # package name
  version: origin/develop                        # tag, commit, or branch
  repo:    https://github.com/Cepave/common.git  # (optional) git URL
```

# Package Release

	make clean all pack
