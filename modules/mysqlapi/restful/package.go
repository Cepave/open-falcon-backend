package restful

import (
	"github.com/gin-gonic/gin"

	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	ov "github.com/Cepave/open-falcon-backend/common/validate"

	graphRest "github.com/Cepave/open-falcon-backend/modules/mysqlapi/restful/graph"
	owlRest "github.com/Cepave/open-falcon-backend/modules/mysqlapi/restful/owl"
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
	mvcBuilder := mvc.NewMvcBuilder(initGinMvcConfig())

	graphRest.InitHttpServices(router, mvcBuilder)

	h := mvcBuilder.BuildHandler

	v1 := router.Group("/api/v1")

	v1.GET("/metrics/builtin", h(getBuiltinMetrics))
	v1.GET("/strategies", h(getStrategies))
	v1.GET("/expressions", h(getExpressions))

	v1.GET("/nqm/agents", h(listAgents))
	v1.GET("/nqm/agent/:agent_id", getAgentById)
	v1.POST("/heartbeat/nqm/agent", h(nqmAgentHeartbeat))
	v1.GET("/heartbeat/nqm/agent/:agent_id/targets", h(nqmAgentHeartbeatTargetList))
	v1.POST("/nqm/agent", addNewAgent)
	v1.PUT("/nqm/agent/:agent_id", modifyAgent)
	v1.POST("/nqm/agent/:agent_id/pingtask", addPingtaskToAgentForAgent)
	v1.DELETE("/nqm/agent/:agent_id/pingtask/:pingtask_id", removePingtaskFromAgentForAgent)
	v1.GET("/nqm/agent/:agent_id/targets", h(listTargetsOfAgentById))
	v1.POST("/nqm/agent/:agent_id/targets/clear", h(clearCachedTargetsOfAgentById))

	v1.GET("/nqm/pingtasks", h(listPingtasks))
	v1.GET("/nqm/pingtask/:pingtask_id", h(getPingtasksById))
	v1.POST("/nqm/pingtask", h(addNewPingtask))
	v1.PUT("/nqm/pingtask/:pingtask_id", h(modifyPingtask))
	v1.POST("/nqm/pingtask/:pingtask_id/agent", addPingtaskToAgentForPingtask)
	v1.DELETE("/nqm/pingtask/:pingtask_id/agent/:agent_id", removePingtaskFromAgentForPingtask)

	v1.GET("/nqm/targets", h(listTargets))
	v1.GET("/nqm/target/:target_id", getTargetById)
	v1.POST("/nqm/target", addNewTarget)
	v1.PUT("/nqm/target/:target_id", modifyTarget)

	v1.GET("/nqm/pingtask/:pingtask_id/agents", h(listAgentsByPingTask))

	v1.GET("/owl/isps", listISPs)
	v1.GET("/owl/isp/:isp_id", h(getISPByID))
	v1.GET("/owl/provinces", listProvinces)
	v1.GET("/owl/province/:province_id", h(getProvinceByID))
	v1.GET("/owl/cities", listCities)
	v1.GET("/owl/city/:city_id", h(getCityByID))
	v1.GET("/owl/province/:province_id/cities", listCitiesInProvince)

	v1.GET("/owl/nametags", h(listNameTags))
	v1.GET("/owl/nametag/:name_tag_id", h(getNameTagById))

	v1.GET("/owl/grouptags", h(listGroupTags))
	v1.GET("/owl/grouptag/:group_tag_id", h(getGroupTagById))

	v1.GET("/hosts", h(listHosts))
	v1.GET("/hostgroups", h(listHostgroups))
	v1.GET("/agent/config", h(getAgentConfig))
	v1.GET("/agent/plugins/:agent_hostname", h(getPlugins))
	v1.GET("/agent/mineplugins", h(getMinePlugins))
	v1.POST("/agent/heartbeat", h(falconAgentHeartbeat))

	v1.GET("/owl/query-object/:uuid", h(owlRest.GetQueryObjectByUuid))
	v1.POST("/owl/query-object", h(owlRest.SaveQueryObject))
	v1.POST("/owl/query-object/vacuum", h(owlRest.VacuumOldQueryObjects))

	router.GET("/health", h(health))
}

func initGinMvcConfig() *mvc.MvcConfig {
	newConfig := mvc.NewDefaultMvcConfig()
	ov.RegisterDefaultValidators(newConfig.Validator)

	return newConfig
}
