package imdb

import (
	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/imdb"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func GetStrValue(c *gin.Context) {
	id, err := httpparams.GetInt(c.Params, "id")
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	value := imdb.StrValue{ID: *id}

	if err := db.IMDB.Find(&value); err.Error != nil {
		h.JSONR(c, badstatus, err.Error)
		return
	}
	h.JSONR(c, h.DataWaper{Data: value})
}

type UpdateStrValueInput struct {
	Value string `json:"value" form:"value" binding:"required"`
}

func UpdateStrValue(c *gin.Context) {
	id, err := httpparams.GetInt(c.Params, "id")
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	var inputs UpdateStrValueInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	value := imdb.StrValue{ID: *id}

	if err := db.IMDB.Find(&value); err.Error != nil {
		h.JSONR(c, badstatus, err.Error)
		return
	}
	copier.Copy(&value, &inputs)
	if err := db.IMDB.Model(&value).Update(&value); err.Error != nil {
		h.JSONR(c, badstatus, err.Error)
		return
	}
	h.JSONR(c, h.DataWaper{Data: value})
}
