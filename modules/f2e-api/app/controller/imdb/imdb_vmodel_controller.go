package imdb

import (
	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/imdb"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	mypg "github.com/masato25/mygo_pagination"
	log "github.com/sirupsen/logrus"
)

func ValueModels(c *gin.Context) {
	var values []imdb.Tag
	if err := db.IMDB.Where("tag_type_id = ?", 3).Preload("TagType").Find(&values); err.Error != nil {
		h.JSONR(c, badstatus, err.Error)
		return
	}
	h.JSONR(c, h.DataWaper{Data: values})
}

type GetValueModelInput struct {
	h.FilterBy
	mypg.Pagging
}

func GetValueModel(c *gin.Context) {
	id, err := httpparams.GetInt(c.Params, "id")
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	inputs := GetValueModelInput{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	var values []imdb.ValueModel
	dt := db.IMDB.Model(&values).Where("tag_id = ?", id)
	dt = inputs.FilterLike(dt, "value")
	pg, _ := inputs.GenOffset(dt)
	if err := dt.Offset(pg.Offset).Limit(pg.Limit).Find(&values).Error; err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	h.JSONR(c, h.DataWaper{Data: values})
}

type CreateValueModelInput struct {
	Value string `json:"value" form:"value" binding:"required"`
}

func CreateValueModel(c *gin.Context) {
	id, err := httpparams.GetInt(c.Params, "id")
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	var inputs CreateValueModelInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	valuem := imdb.ValueModel{TagId: *id}
	copier.Copy(&valuem, &inputs)
	if err := db.IMDB.Save(&valuem); err.Error != nil {
		h.JSONR(c, badstatus, err.Error)
		return
	}
	h.JSONR(c, h.DataWaper{
		Msg:  "ok",
		Data: valuem,
	})
}

type UpdateValueModelInput struct {
	ID    int    `json:"id" form:"id" binding:"required"`
	Value string `json:"value" form:"value" binding:"required"`
}

func UpdateValueModel(c *gin.Context) {
	id, err := httpparams.GetInt(c.Params, "id")
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	var inputs UpdateValueModelInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	valuem := imdb.ValueModel{ID: inputs.ID}
	if err := db.IMDB.Model(&valuem).Where("tag_id = ?", id).Find(&valuem); err.Error != nil {
		h.JSONR(c, badstatus, err.Error)
		return
	}

	copier.Copy(&valuem, &inputs)
	if err := db.IMDB.Model(&valuem).Update(&valuem); err.Error != nil {
		h.JSONR(c, badstatus, err.Error)
		return
	}
	h.JSONR(c, h.DataWaper{
		Msg:  "ok",
		Data: valuem,
	})
}

type DeleteValueModelInput struct {
	ID []int `json:"value_model_ids" form:"value_model_ids" binding:"required"`
	// test only
	Test bool `json:"test"`
}

func DeleteValueModel(c *gin.Context) {
	id, err := httpparams.GetInt(c.Params, "id")
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	var inputs DeleteValueModelInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	log.Debugf("id: %v, inputs.ID: %v", id, inputs.ID)
	var values []imdb.ValueModel
	dt := db.IMDB.Begin()
	if err := dt.Model(&values).Where("tag_id = ? and id in (?)", id, inputs.ID).Find(&values); err.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, "got error when find value model: "+err.Error.Error())
		return
	}
	var valId = make([]int, len(values))
	for indx, vt := range values {
		valId[indx] = vt.ID
	}
	var vmodelVal []imdb.VmodelValue
	objTagsDeletedNum := int64(0)
	dt.Where("value_model_id in (?)", valId).Find(&vmodelVal)
	// if vmodel_values is not empty will delete related object_tags first
	if len(vmodelVal) != 0 {
		valId = make([]int, len(vmodelVal))
		for indx, vt := range vmodelVal {
			valId[indx] = vt.ID
		}
		objTags := []imdb.ObjectTag{}
		dt2 := dt.Model(&objTags).Where("tag_id = ? and value_id in (?)", vmodelVal[0].TagId, valId).Delete(&objTags)
		if dt.Error != nil {
			dt.Rollback()
			h.JSONR(c, badstatus, "got error when delete object_tags: "+dt2.Error.Error())
			return
		}
		objTagsDeletedNum = dt2.RowsAffected
	}
	dt2 := dt.Where("tag_id = ? and id in (?)", id, inputs.ID).Delete(&values)
	if dt2.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, "got error when delete value model: "+dt2.Error.Error())
		return
	}
	if !inputs.Test {
		dt.Commit()
	}
	h.JSONR(c, h.DataWaper{
		Msg: "ok",
		Data: map[string]interface{}{
			"deleted_number_of_object_tags":  objTagsDeletedNum,
			"deleted_number_of_value_models": dt2.RowsAffected,
		},
	})
}
