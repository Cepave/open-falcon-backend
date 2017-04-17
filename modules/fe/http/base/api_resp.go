package base

import (
	"errors"
	"regexp"

	"github.com/Cepave/open-falcon-backend/modules/fe/model/uic"
)

type ApiResp struct {
	Version string                 `json:"version,omitempty"`
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
		msg.Status = "failed"
	} else {
		msg.Status = "success"
	}

	this.Data["json"] = msg
	this.ServeJSON()
}

func (this *BaseController) SessionCheck() (session *uic.Session, err error) {
	name := this.GetString("cName", this.Ctx.GetCookie("name"))
	sig := this.GetString("cSig", this.Ctx.GetCookie("sig"))
	if sig == "" || name == "" {
		err = errors.New("name or sig is empty, please check again")
		return
	}
	session = uic.ReadSessionBySig(sig)
	if session == nil {
		err = errors.New("can not find this kind of session")
		return
	}
	u := uic.SelectUserById(session.Uid)
	switch {
	case sig == "disabled":
		err = errors.New("this account has beend disabled")
		return
	case u != nil && u.Role == -1:
		err = errors.New("this account has beend disabled")
		return
	case u == nil || u.Name != name:
		err = errors.New("can not find this kind of session")
		return
	default:
		return
	}
}

func (this *BaseController) ResposeError(apiBasicParams *ApiResp, msg string) {
	apiBasicParams.Error["message"] = msg
	this.ServeApiJson(apiBasicParams)
}
