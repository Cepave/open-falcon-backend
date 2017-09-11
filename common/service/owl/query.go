package owl

import (
	"net/http"
	"time"

	"github.com/juju/errors"
	"github.com/satori/go.uuid"
	gt "gopkg.in/h2non/gentleman.v2"

	"github.com/Cepave/open-falcon-backend/common/db"
	oHttp "github.com/Cepave/open-falcon-backend/common/http"
	"github.com/Cepave/open-falcon-backend/common/http/client"
	ojson "github.com/Cepave/open-falcon-backend/common/json"
	model "github.com/Cepave/open-falcon-backend/common/model/owl"
)

type QueryServiceConfig struct {
	*oHttp.RestfulClientConfig
}

type QueryService interface {
	LoadQueryByUuid(uuid uuid.UUID) *model.Query
	CreateOrLoadQuery(query *model.Query)
}

func NewQueryService(config QueryServiceConfig) QueryService {
	newClient := oHttp.NewApiService(config.RestfulClientConfig).NewClient()

	return &queryServiceImpl{
		loadQueryByUuid:   newClient.Get().AddPath("/api/v1/owl/query-object"),
		createOrLoadQuery: newClient.Post().AddPath("/api/v1/owl/query-object"),
	}
}

type queryServiceImpl struct {
	loadQueryByUuid   *gt.Request
	createOrLoadQuery *gt.Request
}

func (s *queryServiceImpl) LoadQueryByUuid(uuid uuid.UUID) *model.Query {
	req := s.loadQueryByUuid.Clone().
		AddPath("/" + uuid.String())

	resp := client.ToGentlemanReq(req).SendAndMustMatch(
		func(resp *gt.Response) error {
			switch resp.StatusCode {
			case http.StatusOK, http.StatusNotFound:
				return nil
			}

			return errors.Errorf(client.ToGentlemanResp(resp).ToDetailString())
		},
	)

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	jsonQuery := &struct {
		Uuid    ojson.Uuid `json:"uuid"`
		NamedId string     `json:"feature_name"`

		Content    ojson.VarBytes `json:"content"`
		Md5Content ojson.Bytes16  `json:"md5_content"`

		CreationTime ojson.JsonTime `json:"creation_time"`
		AccessTime   ojson.JsonTime `json:"access_time"`
	}{}

	client.ToGentlemanResp(resp).MustBindJson(jsonQuery)

	return &model.Query{
		Uuid:    db.DbUuid(jsonQuery.Uuid),
		NamedId: jsonQuery.NamedId,

		Content:    []byte(jsonQuery.Content),
		Md5Content: db.Bytes16(jsonQuery.Md5Content),

		CreationTime: time.Time(jsonQuery.CreationTime),
		AccessTime:   time.Time(jsonQuery.AccessTime),
	}
}

// Loads object of query or creating one.
//
// Any error would be expressed by panic.
func (s *queryServiceImpl) CreateOrLoadQuery(query *model.Query) {
	req := s.createOrLoadQuery.Clone().
		JSON(map[string]interface{}{
			"feature_name": query.NamedId,
			"content":      ojson.VarBytes(query.Content),
			"md5_content":  ojson.Bytes16(query.Md5Content),
		})

	resp := client.ToGentlemanReq(req).SendAndStatusMustMatch(http.StatusOK)

	jsonBody := &struct {
		Uuid         ojson.Uuid     `json:"uuid"`
		CreationTime ojson.JsonTime `json:"creation_time"`
		AccessTime   ojson.JsonTime `json:"access_time"`
	}{}

	client.ToGentlemanResp(resp).MustBindJson(jsonBody)

	query.Uuid = db.DbUuid(jsonBody.Uuid)
	query.CreationTime = time.Time(jsonBody.CreationTime)
	query.AccessTime = time.Time(jsonBody.AccessTime)
}
