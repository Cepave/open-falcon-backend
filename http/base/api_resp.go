package base

type ApiResp struct {
	Version string                 `json:"value,omitempty"`
	Method  string                 `json:"method,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Error   map[string]interface{} `json:"error,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func (this *BaseController) ApiBasicParams() (apiResp *ApiResp, errorTmp, dataTmp map[string]interface{}) {
	apiResp = new(ApiResp)
	apiResp.Version = "v1"
	errorTmp = map[string]interface{}{}
	dataTmp = map[string]interface{}{}
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
