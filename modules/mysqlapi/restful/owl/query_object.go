package owl

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"

	"github.com/Cepave/open-falcon-backend/common/db"
	ogin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
	"github.com/Cepave/open-falcon-backend/common/json"
	"github.com/Cepave/open-falcon-backend/common/model/owl"

	owlDb "github.com/Cepave/open-falcon-backend/modules/mysqlapi/rdb/owl"
	owlSrv "github.com/Cepave/open-falcon-backend/modules/mysqlapi/service/owl"
)

func SaveQueryObject(jsonQuery *QueryObject) mvc.OutputBody {
	queryModel := jsonQuery.toQueryModel()

	owlSrv.QueryObjectService.CreateOrLoadQuery(queryModel)

	return mvc.JsonOutputBody(queryModel.ToJson())
}

func GetQueryObjectByUuid(
	d *struct {
		Uuid uuid.UUID `mvc:"param[uuid]"`
		Uri  string    `mvc:"req[RequestURI]"`
	},
) mvc.OutputBody {
	queryObject := owlSrv.QueryObjectService.LoadQueryByUuid(d.Uuid)

	if queryObject == nil {
		uuidString := d.Uuid.String()
		return mvc.JsonOutputBody2(
			http.StatusNotFound,
			map[string]interface{}{
				"uuid":        uuidString,
				"http_status": http.StatusNotFound,
				"error_code":  1,
				"uri":         d.Uri,
			},
		)
	}

	return mvc.JsonOutputBody(queryObject.ToJson())
}

func VacuumOldQueryObjects(
	d *struct {
		ForDays int `mvc:"query[for_days] default[180]" validate:"min=14"`
	},
) mvc.OutputBody {
	beforeTime := time.Now().Add(time.Duration(-d.ForDays) * time.Duration(24) * time.Hour)

	affectedRows := owlDb.RemoveOldQueryObject(beforeTime)

	return mvc.JsonOutputBody(
		map[string]interface{}{
			"before_time":   beforeTime.Unix(),
			"affected_rows": affectedRows,
		},
	)
}

type QueryObject struct {
	NamedId string `json:"feature_name" validate:"min=1"`

	Content    json.VarBytes `json:"content" validate:"min=1"`
	Md5Content json.Bytes16  `json:"md5_content" validate:"non_zero_slice"`
}

func (qb *QueryObject) Bind(c *gin.Context) {
	ogin.BindJson(c, qb)
}

func (qb *QueryObject) toQueryModel() *owl.Query {
	return &owl.Query{
		NamedId:    qb.NamedId,
		Content:    []byte(qb.Content),
		Md5Content: db.Bytes16(qb.Md5Content),
	}
}
