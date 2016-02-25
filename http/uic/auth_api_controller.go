package uic

import (
	"errors"
	"github.com/Cepave/fe/g"
	"github.com/Cepave/fe/http/base"
	. "github.com/Cepave/fe/model/uic"
	"github.com/Cepave/fe/utils"
	"github.com/toolkits/str"
	"strings"
	"time"
)

type AuthApiController struct {
	base.BaseController
}

func (this *AuthApiController) ResposeError(apiBasicParams *base.ApiResp, msg string) {
	apiBasicParams.Error["message"] = msg
	this.ServeApiJson(apiBasicParams)
}

func (this *AuthApiController) SessionCheck(name, sig string) (session *Session, err error) {
	if sig == "" || name == "" {
		err = errors.New("name or sig is empty, please check again")
		return
	}
	session = ReadSessionBySig(sig)
	if session.Uid != SelectUserIdByName(name) {
		err = errors.New("can not find this kind of session")
		return
	}
	return
}

func (this *AuthApiController) AuthSession() {
	baseResp := this.BasicRespGen()
	name := this.GetString("name", this.Ctx.GetCookie("name"))
	sig := this.GetString("sig", this.Ctx.GetCookie("sig"))
	session, err := this.SessionCheck(name, sig)
	switch {
	case err != nil:
		this.ResposeError(baseResp, err.Error())
	case session.Sig != "":
		baseResp.Data["sigs"] = session
		baseResp.Data["expired"] = session.Expired
		baseResp.Data["message"] = "session passed!"
	default:
		baseResp.Error["message"] = "sesion checking failed for a unknow reason, please ask administor for help"
	}
	this.ServeApiJson(baseResp)
}

func (this *AuthApiController) Logout() {
	baseResp := this.BasicRespGen()
	name := this.GetString("name", this.Ctx.GetCookie("name"))
	sig := this.GetString("sig", this.Ctx.GetCookie("sig"))
	session, err := this.SessionCheck(name, sig)
	switch {
	case err != nil:
		this.ResposeError(baseResp, err.Error())
	default:
		_, err := RemoveSessionByUid(session.Uid)
		if err != nil {
			this.ResposeError(baseResp, err.Error())
		} else {
			baseResp.Data["message"] = "session is removed"
		}
	}
	this.ServeApiJson(baseResp)
}

func (this *AuthApiController) Login() {
	baseResp := this.BasicRespGen()
	name := this.GetString("name", "")
	password := this.GetString("password", "")

	if name == "" || password == "" {
		this.ResposeError(baseResp, "name or password is blank")
	}

	user := ReadUserByName(name)
	switch {
	case user == nil:
		this.ResposeError(baseResp, "no such user")
	case user.Passwd != str.Md5Encode(g.Config().Salt+password):
		this.ResposeError(baseResp, "password error")
	}

	appSig := this.GetString("sig", "")
	callback := this.GetString("callback", "")
	sig, expired := ReadSessionByName(name)
	switch {
	case sig != "":
		baseResp.Data["sig"] = sig
		baseResp.Data["expired"] = expired
	case appSig != "" && callback != "":
		SaveSessionAttrs(user.Id, appSig, int(time.Now().Unix())+3600*24*30)
		baseResp.Data["sig"] = appSig
		baseResp.Data["expired"] = int(time.Now().Unix()) + 3600*24*30
	default:
		sig, expired := this.CreateSession(user.Id, 3600*24*30)
		baseResp.Data["sig"] = sig
		baseResp.Data["expired"] = expired
	}
	this.ServeApiJson(baseResp)
}

func (this *AuthApiController) Register() {
	baseResp := this.BasicRespGen()
	if !g.Config().CanRegister {
		this.ResposeError(baseResp, "registration system is not open")
		return
	}

	name := strings.TrimSpace(this.GetString("name", ""))
	email := strings.TrimSpace(this.GetString("email", ""))
	password := strings.TrimSpace(this.GetString("password", ""))
	repeatPassword := strings.TrimSpace(this.GetString("repeat_password", ""))

	var lastID int64
	var err error
	switch {
	case password != repeatPassword:
		this.ResposeError(baseResp, "password not equal the repeart one")
	case !utils.IsUsernameValid(name):
		this.ResposeError(baseResp, "name pattern is invalid")
	case ReadUserIdByName(name) > 0:
		this.ResposeError(baseResp, "name is already existent")
	default:
		lastID, err = InsertRegisterUser(name, str.Md5Encode(g.Config().Salt+password), email)
		if err != nil {
			this.ResposeError(baseResp, "insert user fail "+err.Error())
		}
	}

	sig, expired := this.CreateSession(lastID, 3600*24*30)
	baseResp.Data["sig"] = sig
	baseResp.Data["expired"] = expired
	this.ServeApiJson(baseResp)
}

func (this *AuthApiController) CreateSession(uid int64, maxAge int) (sig string, expired int) {
	sig = utils.GenerateUUID()
	expired = int(time.Now().Unix()) + maxAge
	SaveSessionAttrs(uid, sig, expired)
	return
}
