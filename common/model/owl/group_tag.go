package owl

import (
	"fmt"
)

type GroupTag struct {
	Id int32 `gorm:"primary_key:true;column:gt_id"`
	Name string `gorm:"column:gt_name"`
}

func (GroupTag) TableName() string {
	return "owl_group_tag"
}

func (groupTag *GroupTag) String() string {
	return fmt.Sprintf("Group Tag[%d][%s]", groupTag.Id, groupTag.Name)
}
