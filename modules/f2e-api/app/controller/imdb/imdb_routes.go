package imdb

import (
	"net/http"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	"github.com/gin-gonic/gin"
)

var db config.DBPool
var httpparams = h.Params{}

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = config.Con()
	imdbr := r.Group("/api/v1/imdb")
	// tags
	imdbr.GET("/tag_types", GetTagType)
	imdbr.GET("/tags", TagLists)
	imdbr.GET("/tag/:id", GetTag)
	imdbr.POST("/tag", CreateTag)
	imdbr.PUT("/tag/:id", UpdateTag)
	imdbr.DELETE("/tag/:id", DeleteTag)

	// vmodels
	imdbr.GET("/value_models", ValueModels)
	imdbr.GET("/value_model/:id", GetValueModel)
	imdbr.POST("/value_model/:id", CreateValueModel)
	imdbr.PUT("/value_model/:id", UpdateValueModel)
	imdbr.DELETE("/value_models/:id", DeleteValueModel)

	// str_value
	// imdbr.GET("/str_value/:id", GetStrValue)
	// imdbr.POST("/str_value", CreateStrValue)
	// imdbr.PUT("/str_value", UpdateStrValue)

	// imdbr.GET("/int_value/:id", GetIntValue)
	// imdbr.POST("/int_value", CreateIntValue)
	// imdbr.PUT("/int_value", UpdateIntValue)

	// imdbr.GET("/description_value/:id", GetDescriptionValue)
	// imdbr.POST("/description_value", CreateDescriptionValue)
	// imdbr.PUT("/description_value", UpdateDescriptionValue)

	// imdbr.GET("/vmodel_value/:id", GetVmodelValue)
	// imdbr.POST("/vmodel_value", CreateVmodelValue)
	// imdbr.PUT("/vmodel_value", UpdateVmodelValue)

	imdbr.GET("/object_tags", ObjectTagList)
	imdbr.GET("/object_tag/:id", GetObjectTag)
	imdbr.POST("/object_tag", CreateObjectTag)
	imdbr.PUT("/object_tag/:id", UpdateObjectTag)
	imdbr.DELETE("/object_tag/:id", DeleteObjctTag)

	// imdbr.POST("/object_tags", APISync)

	// imdbr.POST("/query/object_tag")
	// imdbr.POST("/query/dsl_query")

	// imdbr.GET("resource_objects")
	// imdbr.GET("resource_object/:id")

	// imdbr.GET("dump")

	// imdbr.GET("host_groups")
	// imdbr.POST("host_group")
	// imdbr.PUT("host_group/:id")
}
