package owl

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/model"
)

// Use the pre-defined column name to generate syntax of "ORDER BY XX"
//
// 	gt_number - The number of group tags
// 	gt_names - The ordered name of group tags
func GetSyntaxOfOrderByGroupTags(entity *model.OrderByEntity) string {
	var dirOfGroupTags = "DESC"
	if entity.Direction == model.Ascending {
		dirOfGroupTags = "ASC"
	}

	return fmt.Sprintf("gt_number %s, gt_names %s", dirOfGroupTags, dirOfGroupTags)
}
