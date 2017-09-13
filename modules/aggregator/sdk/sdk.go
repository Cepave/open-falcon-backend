package sdk

import (
	"encoding/json"
	"fmt"
	"time"

	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/common/sdk/requests"
	"github.com/Cepave/open-falcon-backend/modules/aggregator/g"
	f "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/falcon_portal"
	"github.com/toolkits/net/httplib"
)

func HostnamesByID(group_id int64) ([]string, error) {

	uri := fmt.Sprintf("%s/api/v1/hostgroup/%d", g.Config().Api.PlusApi, group_id)
	req, err := requests.CurlPlus(uri, "GET", "aggregator", g.Config().Api.PlusApiToken,
		map[string]string{}, map[string]string{})

	if err != nil {
		return []string{}, err
	}

	type RESP struct {
		HostGroup f.HostGroup `json:"hostgroup"`
		Hosts     []f.Host    `json:"hosts"`
	}

	resp := &RESP{}
	err = req.ToJson(&resp)
	if err != nil {
		return []string{}, err
	}

	hosts := []string{}
	for _, x := range resp.Hosts {
		hosts = append(hosts, x.Hostname)
	}
	return hosts, nil
}

func QueryLastPoints(endpoints, counters []string) (resp []*cmodel.GraphLastResp, err error) {
	cfg := g.Config()
	uri := fmt.Sprintf("%s/api/v1/graph/graph/lastpoint", cfg.Api.PlusApi)

	var req *httplib.BeegoHttpRequest
	headers := map[string]string{"Content-type": "application/json"}
	req, err = requests.CurlPlus(uri, "POST", "aggregator", cfg.Api.PlusApiToken,
		headers, map[string]string{})

	if err != nil {
		return
	}

	req.SetTimeout(time.Duration(cfg.Api.ConnectTimeout)*time.Millisecond,
		time.Duration(cfg.Api.RequestTimeout)*time.Millisecond)

	body := []*cmodel.GraphLastParam{}
	for _, e := range endpoints {
		for _, c := range counters {
			body = append(body, &cmodel.GraphLastParam{e, c})
		}
	}

	b, err := json.Marshal(body)
	if err != nil {
		return
	}
	req.Body(b)

	err = req.ToJson(&resp)
	if err != nil {
		return
	}

	return resp, nil
}
