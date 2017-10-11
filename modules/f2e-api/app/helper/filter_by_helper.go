package helper

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type FilterBy struct {
	Q string `json:"q" form:"q"`
}

func (self FilterBy) FilterRegexp(q *gorm.DB, feildName string) (qt *gorm.DB) {
	if self.Q == "" {
		qt = q
		return
	}
	qt = q.Where(fmt.Sprintf(" %s regexp ?", feildName), self.Q)
	return
}

func (self FilterBy) FilterLike(q *gorm.DB, feildName string) (qt *gorm.DB) {
	if self.Q == "" {
		qt = q
		return
	}
	qt = q.Where(fmt.Sprintf(" %s like ?", feildName), "%"+self.Q+"%")
	return
}
