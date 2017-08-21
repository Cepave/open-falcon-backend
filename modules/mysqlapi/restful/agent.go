package restful

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	commonNqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	mvc "github.com/Cepave/open-falcon-backend/common/gin/mvc"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	commonNqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	"github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/service/hbscache"
)

func addNewAgent(c *gin.Context) {
	/**
	 * Builds data from body of request
	 */
	agentForAdding := commonNqmModel.NewAgentForAdding()

	commonGin.BindJson(c, agentForAdding)
	commonGin.ConformAndValidateStruct(agentForAdding, commonNqmModel.Validator)
	agentForAdding.UniqueGroupTags()
	// :~)

	newAgent, err := commonNqmDb.AddAgent(agentForAdding)
	if err != nil {
		switch err.(type) {
		case commonNqmDb.ErrDuplicatedNqmAgent:
			commonGin.JsonConflictHandler(
				c,
				commonGin.DataConflictError{
					ErrorCode:    1,
					ErrorMessage: err.Error(),
				},
			)
		default:
			panic(err)
		}

		return
	}

	c.JSON(http.StatusOK, newAgent)
}

func modifyAgent(c *gin.Context) {
	/**
	 * Loads agent from database
	 */
	agentId, agentIdErr := strconv.Atoi(c.Param("agent_id"))
	if agentIdErr != nil {
		panic(agentIdErr)
	}

	originalAgent := commonNqmDb.GetAgentById(int32(agentId))
	if originalAgent == nil {
		commonGin.JsonNoMethodHandler(c)
		return
	}
	// :~)

	/**
	 * Binding JSON body to modified agent
	 */
	modifiedAgent := originalAgent.ToAgentForAdding()
	commonGin.BindJson(c, modifiedAgent)
	commonGin.ConformAndValidateStruct(modifiedAgent, commonNqmModel.Validator)
	modifiedAgent.UniqueGroupTags()
	// :~)

	updatedAgent, err := commonNqmDb.UpdateAgent(originalAgent, modifiedAgent)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, updatedAgent)
}

func getAgentById(c *gin.Context) {
	agentId, err := strconv.Atoi(c.Param("agent_id"))
	if err != nil {
		commonGin.OutputJsonIfNotNil(c, nil)
	}

	agent := commonNqmDb.GetAgentById(int32(agentId))

	commonGin.OutputJsonIfNotNil(c, agent)
}

func listAgents(
	query *commonNqmModel.AgentQuery,
	paging *struct {
		Page *commonModel.Paging `mvc:"pageSize[50] pageOrderBy[status#desc:connection_id#asc]"`
	},
) (*commonModel.Paging, mvc.OutputBody) {
	agents, resultPaging := commonNqmDb.ListAgents(query, *paging.Page)

	return resultPaging, mvc.JsonOutputBody(agents)
}

func listTargetsOfAgentById(
	q *commonNqmModel.TargetsOfAgentQuery,
	p *struct {
		Paging *commonModel.Paging `mvc:"pageSize[50] pageOrderBy[status#desc:name#asc:host#asc]"`
	},
) (*commonModel.Paging, mvc.OutputBody) {
	targetList, resultPaging := commonNqmDb.ListTargetsOfAgentById(q, *p.Paging)
	if targetList != nil {
		targetList.CacheLifeTime = cacheConfig.Lifetime
	}

	return resultPaging, mvc.JsonOutputOrNotFound(targetList)
}

func clearCachedTargetsOfAgentById(
	q *struct {
		AgentID int32 `mvc:"param[agent_id]"`
	},
) mvc.OutputBody {
	r := commonNqmDb.DeleteCachedTargetsOfAgentById(q.AgentID)
	return mvc.JsonOutputOrNotFound(r)
}

func getMinePlugins(
	p *struct {
		Hostname string `mvc:"query[hostname]" validate:"required"`
	},
) mvc.OutputBody {
	reply := &commonModel.NewAgentPluginsResponse{}

	reply.Plugins = hbscache.GetPlugins(p.Hostname)
	reply.Timestamp = time.Now().Unix()
	reply.GitRepo = hbscache.GitRepo.Get()

	return mvc.JsonOutputBody(reply)
}

func getPlugins(
	p *struct {
		Hostname string `mvc:"param[agent_hostname]" validate:"required"`
	},
) mvc.OutputBody {
	plugins := hbscache.GetPlugins(p.Hostname)
	return mvc.JsonOutputBody(plugins)
}

func getBuiltinMetrics(
	p *struct {
		Hostname string `mvc:"query[hostname]" validate:"required"`
		Checksum string `mvc:"query[checksum]"`
	},
) mvc.OutputBody {
	reply := &commonModel.NewBuiltinMetricResponse{}

	metrics, err := hbscache.GetBuiltinMetrics(p.Hostname)
	if err != nil {
		return mvc.JsonOutputBody(nil)
	}

	checksum := ""
	if len(metrics) > 0 {
		checksum = digestBuiltinMetrics(metrics)
	}

	if p.Checksum == checksum {
		reply.Metrics = []*commonModel.NewBuiltinMetric{}
	} else {
		reply.Metrics = metrics
	}
	reply.Checksum = checksum
	reply.Timestamp = time.Now().Unix()

	return mvc.JsonOutputBody(reply)
}

func digestBuiltinMetrics(items []*commonModel.NewBuiltinMetric) string {
	sort.Sort(commonModel.NewBuiltinMetricSlice(items))

	var buf bytes.Buffer
	for _, m := range items {
		buf.WriteString(m.String())
	}

	return utils.Md5(buf.String())
}

func getStrategies() mvc.OutputBody {
	// 一个机器ID对应多个模板ID
	hidTids := hbscache.HostTemplateIds.GetMap()
	sz := len(hidTids)
	if sz == 0 {
		return mvc.JsonOutputBody(nil)
	}

	// Judge需要的是hostname，此处要把HostId转换为hostname
	// 查出的hosts，是不处于维护时间内的
	hosts := hbscache.MonitoredHosts.Get()
	if len(hosts) == 0 {
		// 所有机器都处于维护状态，汗
		return mvc.JsonOutputBody(nil)
	}

	tpls := hbscache.TemplateCache.GetMap()
	if len(tpls) == 0 {
		return mvc.JsonOutputBody(nil)
	}

	strategies := hbscache.Strategies.GetMap()
	if len(strategies) == 0 {
		return mvc.JsonOutputBody(nil)
	}

	// 做个索引，给一个tplId，可以很方便的找到对应了哪些Strategy
	tpl2Strategies := tpl2Strategies(strategies)

	hostStrategies := make([]*commonModel.NewHostStrategy, 0, sz)
	for hostId, tplIds := range hidTids {

		h, exists := hosts[hostId]
		if !exists {
			continue
		}

		// 计算当前host配置了哪些监控策略
		ss := calcInheritStrategies(tpls, tplIds, tpl2Strategies)
		if len(ss) <= 0 {
			continue
		}

		hs := commonModel.NewHostStrategy{
			Hostname:   h.Name,
			Strategies: ss,
		}

		hostStrategies = append(hostStrategies, &hs)

	}

	return mvc.JsonOutputBody(hostStrategies)
}

func tpl2Strategies(strategies map[int]*commonModel.NewStrategy) map[int][]*commonModel.NewStrategy {
	ret := make(map[int][]*commonModel.NewStrategy)
	for _, s := range strategies {
		if s == nil || s.Tpl == nil {
			continue
		}
		if _, exists := ret[s.Tpl.ID]; exists {
			ret[s.Tpl.ID] = append(ret[s.Tpl.ID], s)
		} else {
			ret[s.Tpl.ID] = []*commonModel.NewStrategy{s}
		}
	}
	return ret
}

func calcInheritStrategies(allTpls map[int]*commonModel.NewTemplate, tids []int, tpl2Strategies map[int][]*commonModel.NewStrategy) []*commonModel.NewStrategy {
	// 根据模板的继承关系，找到每个机器对应的模板全量
	/**
	 * host_id =>
	 * |a |d |a |a |a |
	 * |  |  |b |b |f |
	 * |  |  |  |c |  |
	 * |  |  |  |  |  |
	 */
	tpl_buckets := [][]int{}
	for _, tid := range tids {
		ids := hbscache.ParentIds(allTpls, tid)
		if len(ids) <= 0 {
			continue
		}
		tpl_buckets = append(tpl_buckets, ids)
	}

	// 每个host 关联的模板，有继承关系的放到同一个bucket中，其他的放在各自单独的bucket中
	/**
	 * host_id =>
	 * |a |d |a |
	 * |b |  |f |
	 * |c |  |  |
	 * |  |  |  |
	 */
	uniq_tpl_buckets := [][]int{}
	for i := 0; i < len(tpl_buckets); i++ {
		var valid bool = true
		for j := 0; j < len(tpl_buckets); j++ {
			if i == j {
				continue
			}
			if slice_int_eq(tpl_buckets[i], tpl_buckets[j]) {
				break
			}
			if slice_int_lt(tpl_buckets[i], tpl_buckets[j]) {
				valid = false
				break
			}
		}
		if valid {
			uniq_tpl_buckets = append(uniq_tpl_buckets, tpl_buckets[i])
		}
	}

	// 继承覆盖父模板策略，得到每个模板聚合后的策略列表
	strategies := []*commonModel.NewStrategy{}

	exists_by_id := make(map[int]struct{})
	for _, bucket := range uniq_tpl_buckets {

		// 开始计算一个桶，先计算老的tid，再计算新的，所以可以覆盖
		// 该桶最终结果
		bucket_stras_map := make(map[string][]*commonModel.NewStrategy)
		for _, tid := range bucket {

			// 一个tid对应的策略列表
			the_tid_stras := make(map[string][]*commonModel.NewStrategy)

			if stras, ok := tpl2Strategies[tid]; ok {
				for _, s := range stras {
					uuid := fmt.Sprintf("metric:%s/tags:%v", s.Metric, utils.SortedTags(s.Tags))
					if _, ok2 := the_tid_stras[uuid]; ok2 {
						the_tid_stras[uuid] = append(the_tid_stras[uuid], s)
					} else {
						the_tid_stras[uuid] = []*commonModel.NewStrategy{s}
					}
				}
			}

			// 覆盖父模板
			for uuid, ss := range the_tid_stras {
				bucket_stras_map[uuid] = ss
			}
		}

		last_tid := bucket[len(bucket)-1]

		// 替换所有策略的模板为最年轻的模板
		for _, ss := range bucket_stras_map {
			for _, s := range ss {
				valStrategy := s
				// exists_by_id[s.Id] 是根据策略ID去重，不太确定是否真的需要，不过加上肯定没问题
				if _, exist := exists_by_id[valStrategy.ID]; !exist {
					if valStrategy.Tpl.ID != last_tid {
						valStrategy.Tpl = allTpls[last_tid]
					}
					strategies = append(strategies, valStrategy)
					exists_by_id[valStrategy.ID] = struct{}{}
				}
			}
		}
	}

	return strategies
}

func slice_int_contains(list []int, target int) bool {
	for _, b := range list {
		if b == target {
			return true
		}
	}
	return false
}

func slice_int_eq(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, av := range a {
		if av != b[i] {
			return false
		}
	}
	return true
}

func slice_int_lt(a []int, b []int) bool {
	for _, i := range a {
		if !slice_int_contains(b, i) {
			return false
		}
	}
	return true
}

func getExpressions() mvc.OutputBody {
	expressions := hbscache.ExpressionCache.Get()
	return mvc.JsonOutputBody(expressions)
}
