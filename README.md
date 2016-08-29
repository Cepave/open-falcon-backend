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
package: github.com/cepave/open-falcon-backend

import:
- package: github.com/cpeave/common              # package name
  version: origin/develop                        # tag, commit, or branch
  repo:    https://github.com/cepave/common.git  # (optional) git URL
```

# Package Release

	make clean all pack
