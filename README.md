# Documentations

- http://book.open-falcon.org
- http://docs.openfalcon.apiary.io

# Get Started

    git clone https://github.com/cepave/open-falcon-backend.git
    cd open-falcon-backend
    git submodule update --init

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

Make sure you're using Go 1.5+ and **GO15VENDOREXPERIMENT=1** env var is exported. (You can ignore GO15VENDOREXPERIMENT using Go 1.6+.)

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

## 自動偵測模板改變重製告警step
當使用者改變模改內的告警條件 ex. `cpu.idle > 10` -> `cpu.idle > 50`, 告警的step counter將會被重製.

## 新增兩個api
* /api/judges/list
  * 拿到judge內部所有現存告警列表
* /api/judges/delete
  * 使用特定的hash將該告警從judge內部移除
