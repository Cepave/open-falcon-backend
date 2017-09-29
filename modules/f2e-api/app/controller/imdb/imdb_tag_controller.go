package imdb

import (
	"fmt"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/imdb"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	mypg "github.com/masato25/mygo_pagination"
)

func GetTagType(c *gin.Context) {
	var tagType []imdb.TagType
	if err := db.IMDB.Find(&tagType); err != nil {
		h.JSONR(c, err.Error)
		return
	}
	h.JSONR(c, h.DataWaper{Data: tagType})
}

type GetTagListsInput struct {
	h.FilterBy
	h.SortByHelpr
	mypg.Pagging
}

func TagLists(c *gin.Context) {
	inputs := GetTagListsInput{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, err.Error())
		return
	}
	var tags []imdb.Tag
	dt := db.IMDB
	dt = dt.Model(&tags)
	dt = inputs.FilterLike(dt, "name")
	pg, err := inputs.GenOffset(dt)
	if err != nil {
		h.JSONR(c, err.Error())
		return
	}
	dt = inputs.OrderBySql(dt)
	err = dt.Preload("TagType").Limit(pg.Limit).Offset(pg.Offset).Find(&tags).Error
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	h.JSONR(c, h.DataWaper{
		Data: tags,
		Page: pg,
	})
}

type GetTagInput struct {
	ID int `json:"id" form:"id" binding:"required"`
}

func GetTag(c *gin.Context) {
	id, err := httpparams.GetInt(c.Params, "id")
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	tag := imdb.Tag{ID: *id}
	if d := db.IMDB.Preload("TagType").Find(&tag); d.Error != nil {
		h.JSONR(c, badstatus, d.Error)
		return
	}
	h.JSONR(c, h.DataWaper{Data: tag})
}

type CreateTagInput struct {
	Name        string `json:"name" form:"name" binding:"required"`
	TagTypeId   int    `json:"tag_type_id" form:"tag_type_id" binding:"required"`
	Description string `json:"description" from:"description"`
}

func CreateTag(c *gin.Context) {
	var inputs CreateTagInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	tag := imdb.Tag{Default: -1}
	copier.Copy(&tag, &inputs)
	if err := db.IMDB.Save(&tag); err.Error != nil {
		h.JSONR(c, badstatus, err.Error)
		return
	}
	db.IMDB.Preload("TagType").Find(&tag)
	h.JSONR(c, h.DataWaper{
		Msg:  "ok",
		Data: tag,
	})
	return
}

type UpdateTagInput struct {
	ID          int    `json:"id" form:"id" binding:"required"`
	Description string `json:"description" from:"description"`
}

func UpdateTag(c *gin.Context) {
	id, err := httpparams.GetInt(c.Params, "id")
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	inputs := UpdateTagInput{ID: *id}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	tag := imdb.Tag{ID: inputs.ID}
	if err := db.IMDB.Find(&tag); err.Error != nil {
		h.JSONR(c, badstatus, err.Error)
		return
	}
	copier.Copy(&tag, &inputs)
	if err := db.IMDB.Model(&tag).Update(&tag); err.Error != nil {
		h.JSONR(c, badstatus, err.Error)
		return
	}
	db.IMDB.Preload("TagType").Find(&tag)
	h.JSONR(c, h.DataWaper{
		Msg:  "ok",
		Data: tag,
	})
	return
}

func DeleteTag(c *gin.Context) {
	id, err := httpparams.GetInt(c.Params, "id")
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	tag := imdb.Tag{ID: *id}
	if err := db.IMDB.Find(&tag); err.Error != nil {
		h.JSONR(c, badstatus, err.Error)
		return
	}
	if tag.Default == 1 {
		h.JSONR(c, badstatus, fmt.Errorf("tag id: %v is a default tag, can not be remove", tag.ID))
		return
	}
	if err := db.IMDB.Delete(&tag); err.Error != nil {
		h.JSONR(c, badstatus, err.Error)
		return
	}
	h.JSONR(c, h.DataWaper{
		Msg: "ok",
		Data: map[string]int{
			"id": tag.ID,
		},
	})
	return
}
