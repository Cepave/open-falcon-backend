package imdb

import (
	con "github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
)

type ObjectTag struct {
	ID               int      `gorm:"primary_key" json:"id"`
	TagId            int      `gorm:"not null" json:"tag_id"`
	ResourceObjectId int      `gorm:"not null" json:"resource_object_id"`
	ValueId          int      `gorm:"not null" json:"value_id"`
	CreatedBy        string   `json:"created_by"`
	CreatedAt        JSONTime `json:"created_at"`
	UpdatedAt        JSONTime `json:"updated_at"`
}

func (ObjectTag) TableName() string {
	return "object_tags"
}

func (self *ObjectTag) Value() (interface{}, error) {
	db := con.Con()
	tag := Tag{ID: self.TagId}
	if err := db.IMDB.Preload("TagType").Find(&tag); err.Error != nil {
		return nil, err.Error
	}
	switch tag.TagType.TypeName {
	case "value_model":
		mvalue := VmodelValue{ObjectTagId: self.ID}
		if err := db.IMDB.Where("object_tag_id = ?", self.ID).Find(&mvalue); err.Error != nil {
			return nil, err.Error
		}
		vmodel := ValueModel{ID: mvalue.ValueModelId}
		if err := db.IMDB.Find(&vmodel); err.Error != nil {
			return nil, err.Error
		}
		return vmodel.Value, nil
	case "int":
		var res IntValue
		err := db.IMDB.Where("object_tag_id = ?", self.ID).Find(&res)
		return res.Value, err.Error
	case "string":
		var res StrValue
		err := db.IMDB.Where("object_tag_id = ?", self.ID).Find(&res)
		return res.Value, err.Error
	case "description":
		var res DescriptionValue
		err := db.IMDB.Where("object_tag_id = ?", self.ID).Find(&res)
		return res.Value, err.Error
	}
	return nil, nil
}
