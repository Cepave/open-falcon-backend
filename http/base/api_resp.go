package base

import (
	"regexp"
)

type ApiResp struct {
	Version string                 `json:"value,omitempty"`
	Method  string                 `json:"method,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Error   map[string]interface{} `json:"error,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func (this *BaseController) BasicRespGen() (apiResp *ApiResp) {
	apiResp = new(ApiResp)
	r, _ := regexp.Compile("/api/([^/]+)")
	apiResp.Version = r.FindStringSubmatch(this.Ctx.Request.URL.RequestURI())[1]
	apiResp.Method = this.Ctx.Request.Method
	apiResp.Error = map[string]interface{}{}
	apiResp.Data = map[string]interface{}{}
	return
}

func (this *BaseController) ServeApiJson(msg *ApiResp) {
	if len(msg.Error) != 0 {
		msg.Status = "Failed"
	} else {
		msg.Status = "Success"
	}

	this.Data["json"] = msg
	this.ServeJSON()
}
