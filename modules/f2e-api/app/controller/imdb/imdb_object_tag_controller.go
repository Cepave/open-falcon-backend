package imdb

import (
	"errors"
	"fmt"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/imdb"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	mypg "github.com/masato25/mygo_pagination"
)

type ObjectTagsInput struct {
	mypg.Pagging
	h.SortByHelpr
}

func ObjectTagList(c *gin.Context) {
	inputs := ObjectTagsInput{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	var objects []imdb.ObjectTag
	dt := db.IMDB.Model(&imdb.ObjectTag{})
	pg, err := inputs.GenOffset(dt)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	dt = inputs.SortByHelpr.OrderBySql(dt)
	err = dt.Limit(pg.Limit).Offset(pg.Offset).Find(&objects).Error
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	h.JSONR(c, h.DataWaper{
		Data: objects,
		Page: pg,
	})
}

func GetObjectTag(c *gin.Context) {
	id, err := httpparams.GetInt(c.Params, "id")
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	object := imdb.ObjectTag{ID: *id}
	if err := db.IMDB.Find(&object); err.Error != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	objValue, err := object.Value()
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	h.JSONR(c, h.DataWaper{
		Data: map[string]interface{}{
			"object_tag": object,
			"value":      objValue,
		},
	})
}

type CreateObjectTagInput struct {
	CreateObjectTagBase
	ValueInt  *int    `json:"value_int" form:"value_int"`
	ValueText *string `json:"value_text" form:"value_text"`
}

type CreateObjectTagBase struct {
	TagId            int `json:"tag_id" form:"tag_id" binding:"required"`
	ResourceObjectId int `json:"resource_object_id" form:"resource_object_id" binding:"required"`
}

func (self *CreateObjectTagInput) Check(tagType string) (err error) {
	switch tagType {
	case "int", "value_model":
		if self.ValueInt == nil {
			err = errors.New("value_int is required")
		}
	default:
		if self.ValueText == nil {
			err = errors.New("value_text is required.")
		}
	}
	return err
}

type CreateObjectChildHelp struct {
	TagId            int
	ObjectTagId      int
	ResourceObjectId int
}

func CreateObjectTag(c *gin.Context) {
	var inputs CreateObjectTagInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	// TODO CreatedBy
	object := imdb.ObjectTag{
		CreatedBy: "root",
		ValueId:   -1,
	}
	dt := db.IMDB.Begin()
	copier.Copy(&object, &inputs.CreateObjectTagBase)
	// check tag type
	tag := imdb.Tag{ID: object.TagId}
	if dt2 := dt.Preload("TagType").Find(&tag); dt2.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, "got error when finding tag: "+dt2.Error.Error())
		return
	}
	tagType := tag.TagType
	if err := inputs.Check(tagType.TypeName); err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	// store object_tag
	if dt2 := dt.Save(&object); dt2.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, "got error when store objet_tag: "+dt2.Error.Error())
		return
	}
	childStore := CreateObjectChildHelp{
		TagId:            inputs.TagId,
		ObjectTagId:      object.ID,
		ResourceObjectId: object.ResourceObjectId,
	}
	valueId := 0
	var err error
	switch tagType.TypeName {
	case "string":
		valueId, err = storeStrVal(dt, childStore, *inputs.ValueText)
	case "int":
		valueId, err = storeIntVal(dt, childStore, *inputs.ValueInt)
	case "description":
		valueId, err = storeDescVal(dt, childStore, *inputs.ValueText)
	case "value_model":
		valueId, err = storeVmodelVal(dt, childStore, *inputs.ValueInt)
	}
	if err != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, "got error when store child object: "+err.Error())
		return
	}
	object.ValueId = valueId
	if dt2 := dt.Model(&object).Update(&object); dt2.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, "got error when update value_id into object_tag: "+dt2.Error.Error())
		return
	}
	dt.Commit()

	objValue, err := object.Value()
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	h.JSONR(c, h.DataWaper{
		Msg: "ok",
		Data: map[string]interface{}{
			"object_tag": object,
			"value":      objValue,
		},
	})
}

func storeStrVal(dt *gorm.DB, storeVal CreateObjectChildHelp, val string) (id int, err error) {
	strVal := imdb.StrValue{
		ResourceObjectId: storeVal.ResourceObjectId,
		TagId:            storeVal.TagId,
		Value:            val,
		ObjectTagId:      storeVal.ObjectTagId,
	}
	dt2 := dt.Save(&strVal)
	err = dt2.Error
	id = strVal.ID
	return
}

func storeDescVal(dt *gorm.DB, storeVal CreateObjectChildHelp, val string) (id int, err error) {
	desVal := imdb.DescriptionValue{
		ResourceObjectId: storeVal.ResourceObjectId,
		TagId:            storeVal.TagId,
		Value:            val,
		ObjectTagId:      storeVal.ObjectTagId,
	}
	dt2 := dt.Save(&desVal)
	err = dt2.Error
	id = desVal.ID
	return
}

func storeIntVal(dt *gorm.DB, storeVal CreateObjectChildHelp, val int) (id int, err error) {
	intVal := imdb.IntValue{
		ResourceObjectId: storeVal.ResourceObjectId,
		TagId:            storeVal.TagId,
		Value:            val,
		ObjectTagId:      storeVal.ObjectTagId,
	}
	dt2 := dt.Save(&intVal)
	err = dt2.Error
	id = intVal.ID
	return
}

func storeVmodelVal(dt *gorm.DB, storeVal CreateObjectChildHelp, val int) (id int, err error) {
	vmodelVal := imdb.VmodelValue{
		ResourceObjectId: storeVal.ResourceObjectId,
		TagId:            storeVal.TagId,
		ValueModelId:     val,
		ObjectTagId:      storeVal.ObjectTagId,
	}
	dt2 := dt.Save(&vmodelVal)
	err = dt2.Error
	id = vmodelVal.ID
	return
}

type UpdateObjectTagInput struct {
	ValueInt  *int    `json:"value_int" form:"value_int"`
	ValueText *string `json:"value_text" form:"value_text"`
}

func (self *UpdateObjectTagInput) Check(tagType string) (err error) {
	switch tagType {
	case "int", "value_model":
		if self.ValueInt == nil {
			err = errors.New("value_int is required")
		}
	default:
		if self.ValueText == nil {
			err = errors.New("value_text is required.")
		}
	}
	return err
}

func UpdateObjectTag(c *gin.Context) {
	var err error
	//TODO createBy
	inputs := UpdateObjectTagInput{}
	if err = c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	id, err := httpparams.GetInt(c.Params, "id")
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	dt := db.IMDB.Begin()
	object := imdb.ObjectTag{
		ID: *id,
	}
	if dt2 := dt.Find(&object); dt2.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, dt2.Error)
		return
	}

	tag := imdb.Tag{ID: object.TagId}
	if dt2 := dt.Preload("TagType").Find(&tag); dt2.Error != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, dt2.Error)
		return
	}
	err = inputs.Check(tag.TagType.TypeName)
	if err != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, err)
		return
	}

	switch tag.TagType.TypeName {
	case "string":
		err = updateStrVal(dt, object.ValueId, *inputs.ValueText)
	case "int":
		err = updateIntVal(dt, object.ValueId, *inputs.ValueInt)
	case "description":
		err = updateDescVal(dt, object.ValueId, *inputs.ValueText)
	case "value_model":
		err = updateVmodelVal(dt, object.ValueId, *inputs.ValueInt)
	}
	if err != nil {
		dt.Rollback()
		h.JSONR(c, badstatus, err)
		return
	}
	dt.Commit()
	objValue, err := object.Value()
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	h.JSONR(c, h.DataWaper{
		Msg: "ok",
		Data: map[string]interface{}{
			"object_tag": object,
			"value":      objValue,
		},
	})
}

func updateStrVal(dt *gorm.DB, vid int, val string) (err error) {
	strVal := imdb.StrValue{ID: vid}
	dt2 := dt.Find(&strVal)
	if dt2.Error != nil {
		return dt2.Error
	}
	strVal.Value = val
	dt2 = dt.Model(&strVal).Update(&strVal)
	err = dt2.Error
	return
}

func updateDescVal(dt *gorm.DB, vid int, val string) (err error) {
	desVal := imdb.DescriptionValue{ID: vid}
	dt2 := dt.Find(&desVal)
	if dt2.Error != nil {
		return dt2.Error
	}
	desVal.Value = val
	dt2 = dt.Model(&desVal).Update(&desVal)
	err = dt2.Error
	return
}

func updateIntVal(dt *gorm.DB, vid int, val int) (err error) {
	intVal := imdb.IntValue{ID: vid}
	dt2 := dt.Find(&intVal)
	if dt2.Error != nil {
		return dt2.Error
	}
	intVal.Value = val
	dt2 = dt.Model(&intVal).Update(&intVal)
	err = dt2.Error
	return
}

func updateVmodelVal(dt *gorm.DB, vid int, val int) (err error) {
	vmodelVal := imdb.VmodelValue{ID: vid}
	dt2 := dt.Find(&vmodelVal)
	if dt2.Error != nil {
		return dt2.Error
	}
	vmodelVal.ValueModelId = val
	dt2 = dt.Model(&vmodelVal).Update((&vmodelVal))
	err = dt2.Error
	return
}

type DeleteObjctTagInput struct {
	// test flag
	Test bool `json:"test" form:"test"`
}

func DeleteObjctTag(c *gin.Context) {
	id, err := httpparams.GetInt(c.Params, "id")
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	inputs := DeleteObjctTagInput{Test: false}
	c.Bind(&inputs)
	objectTag := imdb.ObjectTag{ID: *id}
	dt := db.IMDB.Begin()
	if err := dt.Delete(&objectTag).Error; err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	if inputs.Test {
		dt.Rollback()
	} else {
		dt.Commit()
	}
	h.JSONR(c, h.DataWaper{
		Msg: fmt.Sprintf("%v has been deleted", *id),
		Data: map[string]interface{}{
			"id": *id,
			"ts": inputs.Test,
		},
	})
}
