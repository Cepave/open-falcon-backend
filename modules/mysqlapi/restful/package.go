package restful

import (
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	"github.com/gin-gonic/gin"
)

var logger = log.NewDefaultLogger("INFO")
var router *gin.Engine = nil
var GinConfig *commonGin.GinConfig = &commonGin.GinConfig{}

func InitGin(config *commonGin.GinConfig) {
	if router != nil {
		return
	}

	router = commonGin.NewDefaultJsonEngine(config)

	logger.Infof("Going to start web service. Listen: %s", config)

	initApi()

	go commonGin.StartServiceOrExit(router, config)

	*GinConfig = *config
}

type CacheConfig struct {
	Size     int
	Lifetime int
}

var cacheConfig *CacheConfig = &CacheConfig{}

func InitCache(config *CacheConfig) {
	cacheConfig = config
}

func initApi() {
	mvcBuilder := mvc.NewMvcBuilder(mvc.NewDefaultMvcConfig())

	v1 := router.Group("/api/v1")

	v1.GET("/metrics/builtin", mvcBuilder.BuildHandler(getBuiltinMetrics))
	v1.GET("/strategies", mvcBuilder.BuildHandler(getStrategies))
	v1.GET("/expressions", mvcBuilder.BuildHandler(getExpressions))

	v1.GET("/nqm/agents", mvcBuilder.BuildHandler(listAgents))
	v1.GET("/nqm/agent/:agent_id", getAgentById)
	v1.POST("/heartbeat/nqm/agent", mvcBuilder.BuildHandler(nqmAgentHeartbeat))
	v1.GET("/heartbeat/nqm/agent/:agent_id/targets", mvcBuilder.BuildHandler(nqmAgentHeartbeatTargetList))
	v1.POST("/nqm/agent", addNewAgent)
	v1.PUT("/nqm/agent/:agent_id", modifyAgent)
	v1.POST("/nqm/agent/:agent_id/pingtask", addPingtaskToAgentForAgent)
	v1.DELETE("/nqm/agent/:agent_id/pingtask/:pingtask_id", removePingtaskFromAgentForAgent)
	v1.GET("/nqm/agent/:agent_id/targets", mvcBuilder.BuildHandler(listTargetsOfAgentById))
	v1.POST("/nqm/agent/:agent_id/targets/clear", mvcBuilder.BuildHandler(clearCachedTargetsOfAgentById))

	v1.GET("/nqm/pingtasks", mvcBuilder.BuildHandler(listPingtasks))
	v1.GET("/nqm/pingtask/:pingtask_id", mvcBuilder.BuildHandler(getPingtasksById))
	v1.POST("/nqm/pingtask", mvcBuilder.BuildHandler(addNewPingtask))
	v1.PUT("/nqm/pingtask/:pingtask_id", mvcBuilder.BuildHandler(modifyPingtask))
	v1.POST("/nqm/pingtask/:pingtask_id/agent", addPingtaskToAgentForPingtask)
	v1.DELETE("/nqm/pingtask/:pingtask_id/agent/:agent_id", removePingtaskFromAgentForPingtask)

	v1.GET("/nqm/targets", mvcBuilder.BuildHandler(listTargets))
	v1.GET("/nqm/target/:target_id", getTargetById)
	v1.POST("/nqm/target", addNewTarget)
	v1.PUT("/nqm/target/:target_id", modifyTarget)

	v1.GET("/nqm/pingtask/:pingtask_id/agents", mvcBuilder.BuildHandler(listAgentsByPingTask))

	v1.GET("/owl/isps", listISPs)
	v1.GET("/owl/isp/:isp_id", mvcBuilder.BuildHandler(getISPByID))
	v1.GET("/owl/provinces", listProvinces)
	v1.GET("/owl/province/:province_id", mvcBuilder.BuildHandler(getProvinceByID))
	v1.GET("/owl/cities", listCities)
	v1.GET("/owl/city/:city_id", mvcBuilder.BuildHandler(getCityByID))
	v1.GET("/owl/province/:province_id/cities", listCitiesInProvince)

	v1.GET("/owl/nametags", mvcBuilder.BuildHandler(listNameTags))
	v1.GET("/owl/nametag/:name_tag_id", mvcBuilder.BuildHandler(getNameTagById))

	v1.GET("/owl/grouptags", mvcBuilder.BuildHandler(listGroupTags))
	v1.GET("/owl/grouptag/:group_tag_id", mvcBuilder.BuildHandler(getGroupTagById))

	v1.GET("/hosts", mvcBuilder.BuildHandler(listHosts))
	v1.GET("/hostgroups", mvcBuilder.BuildHandler(listHostgroups))
	v1.GET("/agent/config", mvcBuilder.BuildHandler(getAgentConfig))
	v1.GET("/agent/mineplugins", mvcBuilder.BuildHandler(getMinePlugins))
	v1.POST("/agent/heartbeat", mvcBuilder.BuildHandler(falconAgentHeartbeat))

	router.GET("/health", mvcBuilder.BuildHandler(health))
}
