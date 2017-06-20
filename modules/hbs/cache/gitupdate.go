package cache

import (
	"sync"

	"github.com/Cepave/open-falcon-backend/modules/hbs/db"
	log "github.com/sirupsen/logrus"
)

type SafeGitRepo struct {
	sync.RWMutex
	gitRepo string
}

var GitRepo = &SafeGitRepo{gitRepo: ""}

func (this *SafeGitRepo) Get() string {
	this.RLock()
	defer this.RUnlock()
	log.Debugln("Read git repo address from cache.GitRepo: ", this.gitRepo)
	return this.gitRepo
}

func (this *SafeGitRepo) Init() {
	cfg, err := db.QueryConfig("git_repo")
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	log.Debugln("Read git repo address from DB: ", cfg.Value)
	this.gitRepo = cfg.Value
}
