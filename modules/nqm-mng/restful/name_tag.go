package restful

import (
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	"github.com/Cepave/open-falcon-backend/common/model"

	dbOwl "github.com/Cepave/open-falcon-backend/common/db/owl"
)

func listNameTags(
	p *struct {
		Value  string        `mvc:"query[value]"`
		Paging *model.Paging `mvc:"pageSize[100] pageOrderBy[value#asc]"`
	},
) (*model.Paging, mvc.OutputBody) {
	return p.Paging,
		mvc.JsonOutputBody(
			dbOwl.ListNameTags(p.Value, p.Paging),
		)
}

func getNameTagById(
	p *struct {
		NameTagId int16 `mvc:"param[name_tag_id]"`
	},
) mvc.OutputBody {
	return mvc.JsonOutputOrNotFound(
		dbOwl.GetNameTagById(p.NameTagId),
	)
}
