package cache

import (
	"strings"
	"sync"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/hbs/db"
	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/sys"
)

func GitRepoUpdateCheck(hostname string) bool {
	if agentUpdateInfo, ok := Agents.Get(hostname); ok {
		hostGitRepo := agentUpdateInfo.ReportRequest.GitRepo
		currGitRepo := GitRepo.Get()
		log.Debugln("host GitRepo of ", hostname, hostGitRepo)
		return (hostGitRepo != currGitRepo)
	}

	return true
}

func GitUpdateCheck(hostname string) bool {
	if agentUpdateInfo, ok := Agents.Get(hostname); ok {
		hostPluginVersion := agentUpdateInfo.ReportRequest.PluginVersion
		newestPluginVersion := pluginHash.Get()
		log.Debugln("hostPluginVersion of ", hostname, hostPluginVersion)
		log.Debugln("newestPluginVersion is:", newestPluginVersion)
		return (hostPluginVersion != newestPluginVersion)
	}

	// need to update
	return true
}

type SafePluginHash struct {
	sync.RWMutex
	hash string
}

func (this *SafePluginHash) Put(commitHash string) {
	this.Lock()
	defer this.Unlock()
	this.hash = commitHash
}

func (this *SafePluginHash) Get() (commitHash string) {
	this.RLock()
	defer this.RUnlock()
	commitHash = this.hash
	return
}

var pluginHash = &SafePluginHash{hash: ""}

type SafeGitRepo struct {
	sync.RWMutex
	gitRepo string
}

func gitLsRemote(gitRepo string, refs string) (string, error) {
	// This function depends on git command
	if resultStr, err := sys.CmdOut("git", "ls-remote", gitRepo, refs); err != nil {
		return "", err
	} else {
		// resultStr should be:
		// cb7a2998571cb25693867afcb24a7331f597768e        refs/heads/master
		strList := strings.Fields(resultStr)
		return strList[0], nil
	}
}

func getNewestPluginHash() {
	for {
		time.Sleep(time.Minute)

		addr := GitRepo.Get()
		log.Debugln("show GitRepo.Get():", addr)
		if !strings.HasPrefix(addr, "http") {
			continue
		}
		if hash, err := gitLsRemote(addr, "refs/heads/master"); err != nil {
			log.Errorln("Error retrieving git-repo:", addr, err)
		} else {
			pluginHash.Put(hash)
			log.Debugln("Get newest plugin hash from: ", addr)
			log.Debugln("Record newest hash as:", hash)
		}
	}
}

var GitRepo = &SafeGitRepo{gitRepo: ""}

func (this *SafeGitRepo) Put(gitRepo string) {
	this.Lock()
	defer this.Unlock()
	this.gitRepo = gitRepo
}

func (this *SafeGitRepo) Get() (gitRepo string) {
	this.RLock()
	defer this.RUnlock()
	gitRepo = this.gitRepo
	return
}

func getGitRepoAddr() {
	for {
		time.Sleep(time.Minute)

		cfg, err := db.QueryConfig("git_repo")
		if err == nil {
			GitRepo.Put(cfg.Value)
			log.Debugln("Read git repo address from DB: ", cfg.Value)
		}
	}
}
