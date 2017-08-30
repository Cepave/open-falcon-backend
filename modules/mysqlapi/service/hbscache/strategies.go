package hbscache

import (
	"strconv"
	"strings"
	"sync"

	"github.com/Cepave/open-falcon-backend/common/model"
	db "github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/hbsdb"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/container/set"
)

type SafeStrategies struct {
	sync.RWMutex
	M map[int]*model.NewStrategy
}

var Strategies = &SafeStrategies{M: make(map[int]*model.NewStrategy)}

func (this *SafeStrategies) GetMap() map[int]*model.NewStrategy {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func (this *SafeStrategies) Init(tpls map[int]*model.NewTemplate) {
	m, err := db.QueryStrategies(tpls)
	if err != nil {
		return
	}

	this.Lock()
	defer this.Unlock()
	this.M = m
}

func GetBuiltinMetrics(hostname string) ([]*model.NewBuiltinMetric, error) {
	ret := []*model.NewBuiltinMetric{}
	hid, exists := HostMap.GetID(hostname)
	if !exists {
		return ret, nil
	}

	gids, exists := HostGroupsMap.GetGroupIds(hid)
	if !exists {
		return ret, nil
	}

	// 根据gids，获取绑定的所有tids
	tidSet := set.NewIntSet()
	for _, gid := range gids {
		tids, exists := GroupTemplates.GetTemplateIds(gid)
		if !exists {
			continue
		}

		for _, tid := range tids {
			tidSet.Add(tid)
		}
	}

	tidSlice := tidSet.ToSlice()
	if len(tidSlice) == 0 {
		return ret, nil
	}

	// 继续寻找这些tid的ParentId
	allTpls := TemplateCache.GetMap()
	for _, tid := range tidSlice {
		pids := ParentIds(allTpls, tid)
		for _, pid := range pids {
			tidSet.Add(pid)
		}
	}

	// 终于得到了最终的tid列表
	tidSlice = tidSet.ToSlice()

	// 把tid列表用逗号拼接在一起
	count := len(tidSlice)
	tidStrArr := make([]string, count)
	for i := 0; i < count; i++ {
		tidStrArr[i] = strconv.Itoa(tidSlice[i])
	}

	return db.QueryBuiltinMetrics(strings.Join(tidStrArr, ","))
}

func ParentIds(allTpls map[int]*model.NewTemplate, tid int) (ret []int) {
	depth := 0
	for {
		if tid <= 0 {
			break
		}

		if t, exists := allTpls[tid]; exists {
			ret = append(ret, tid)
			tid = t.ParentID
		} else {
			break
		}

		depth++
		if depth == 10 {
			log.Println("[ERROR] template inherit cycle. id:", tid)
			return []int{}
		}
	}

	sz := len(ret)
	if sz <= 1 {
		return
	}

	desc := make([]int, sz)
	for i, item := range ret {
		j := sz - i - 1
		desc[j] = item
	}

	return desc
}
