package cache

import (
	"encoding/xml"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/hbs/db"
	log "github.com/Sirupsen/logrus"
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

type xmlEntry struct {
	ID      string `xml:"id"`
	Updated string `xml:"updated"`
}

type xmlData struct {
	EntryList []xmlEntry `xml:"entry"`
}

func getNewestPluginHash() {
	for {
		time.Sleep(time.Minute)

		addr := GitRepo.Get()
		log.Debugln("show GitRepo.Get():", addr)
		if !strings.HasPrefix(addr, "http") {
			continue
		}
		v := xmlData{}
		atomAddr := strings.Replace(addr, ".git", "/commits/master.atom", -1)
		if resp, err := http.Get(atomAddr); err != nil {
			// handle error.
			log.Errorln("Error retrieving resource:", err)
		} else {
			defer resp.Body.Close()
			xml.NewDecoder(resp.Body).Decode(&v)
		}
		if len(v.EntryList) > 0 {
			// update newest Plugin hash
			strlist := strings.Split(v.EntryList[0].ID, "/")
			hash := strlist[len(strlist)-1]
			pluginHash.Put(hash)
			log.Debugln("Get newest plugin hash from atomAddr:", atomAddr)
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
