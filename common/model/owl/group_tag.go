package owl

import (
	"fmt"
	"strings"

	"github.com/Cepave/open-falcon-backend/common/strconv"
	json "github.com/bitly/go-simplejson"
)

type GroupTag struct {
	Id   int32  `gorm:"primary_key:true;column:gt_id" json:"id" db:"gt_id"`
	Name string `gorm:"column:gt_name" json:"name" db:"gt_name"`
}

func (GroupTag) TableName() string {
	return "owl_group_tag"
}

func (groupTag *GroupTag) String() string {
	return fmt.Sprintf("Group Tag[%d][%s]", groupTag.Id, groupTag.Name)
}

func (groupTag *GroupTag) ToJson() *json.Json {
	jsonGroupTag := json.New()
	jsonGroupTag.Set("id", groupTag.Id)
	jsonGroupTag.Set("name", groupTag.Name)

	return jsonGroupTag
}

type GroupTags []*GroupTag

func (groupTags GroupTags) ToJson() []*json.Json {
	jsonGroupTags := make([]*json.Json, len(groupTags))

	for i, groupTag := range groupTags {
		jsonGroupTags[i] = groupTag.ToJson()
	}

	return jsonGroupTags
}
func (groupTags GroupTags) ToNames() []string {
	namesOfGroupTags := make([]string, len(groupTags))
	for i, groupTag := range groupTags {
		namesOfGroupTags[i] = groupTag.Name
	}

	return namesOfGroupTags
}

// Converts a string of ids and a string of names to a array of GroupTags
func SplitToArrayOfGroupTags(
	ids string, splitForIds string,
	names string, splitForNames string,
) []*GroupTag {
	result := make([]*GroupTag, 0)

	if ids == "" {
		return result
	}

	allIds := strconv.SplitStringToIntArray(ids, splitForIds)
	allNames := strings.Split(names, splitForNames)

	for i, groupTagId := range allIds {
		result = append(
			result,
			&GroupTag{
				Id:   int32(groupTagId),
				Name: allNames[i],
			},
		)
	}

	return result
}

type GroupTagOfPingtaskView struct {
	Id   int    `gorm:"primary_key:true;column:gt_id" json:"id" db:"gt_id"`
	Name string `gorm:"column:gt_name" json:"name" db:"gt_name"`
}

func (GroupTagOfPingtaskView) TableName() string {
	return "owl_group_tag"
}
