package restful

import (
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	gin "gopkg.in/gin-gonic/gin.v1"
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

	v1.GET("/nqm/agents", mvcBuilder.BuildHandler(listAgents))
	v1.GET("/nqm/agent/:agent_id", getAgentById)
	v1.POST("/nqm/agent", addNewAgent)
	v1.PUT("/nqm/agent/:agent_id", modifyAgent)
	v1.POST("/nqm/agent/:agent_id/pingtask", addPingtaskToAgentForAgent)
	v1.DELETE("/nqm/agent/:agent_id/pingtask/:pingtask_id", removePingtaskFromAgentForAgent)
	v1.GET("/nqm/agent/:agent_id/targets", mvcBuilder.BuildHandler(listTargetsOfAgentById))

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

	router.GET("/health", health)
}
