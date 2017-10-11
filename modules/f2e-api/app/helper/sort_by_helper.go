package helper

import "github.com/jinzhu/gorm"

type SortByHelpr struct {
	OrderBy string `json:"order_by" form:"order_by"`
}

func (self SortByHelpr) OrderBySql(q *gorm.DB) (qt *gorm.DB) {
	if self.OrderBy == "" {
		qt = q
		return
	}
	qt = q.Order(self.OrderBy)
	return
}
