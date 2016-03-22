package uic

import (
	"strings"
	"time"

	"github.com/Cepave/fe/g"
	"github.com/Cepave/fe/http/base"
	. "github.com/Cepave/fe/model/uic"
	"github.com/Cepave/fe/utils"
	"github.com/toolkits/str"
)

type AuthApiController struct {
	base.BaseController
}

func (this *AuthApiController) AuthSession() {
	baseResp := this.BasicRespGen()
	session, err := this.SessionCheck()
	switch {
	case err != nil:
		this.ResposeError(baseResp, err.Error())
		return
	case session.Sig != "":
		baseResp.Data["sigs"] = session
		baseResp.Data["expired"] = session.Expired
		baseResp.Data["message"] = "session passed!"
	default:
		baseResp.Error["message"] = "sesion checking failed for a unknow reason, please ask administor for help"
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *AuthApiController) Logout() {
	baseResp := this.BasicRespGen()
	session, err := this.SessionCheck()
	switch {
	case err != nil:
		this.ResposeError(baseResp, err.Error())
		return
	default:
		_, err := RemoveSessionByUid(session.Uid)
		if err != nil {
			this.ResposeError(baseResp, err.Error())
			return
		} else {
			baseResp.Data["message"] = "session is removed"
		}
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *AuthApiController) Login() {
	baseResp := this.BasicRespGen()
	name := this.GetString("name", "")
	password := this.GetString("password", "")

	if name == "" || password == "" {
		this.ResposeError(baseResp, "name or password is blank")
		return
	}

	user := ReadUserByName(name)
	switch {
	case user == nil:
		this.ResposeError(baseResp, "no such user")
		return
	case user.Passwd != str.Md5Encode(g.Config().Salt+password):
		this.ResposeError(baseResp, "password error")
		return
	}

	appSig := this.GetString("sig", "")
	callback := this.GetString("callback", "")
	sig, expired := ReadSessionByName(name)
	switch {
	case sig != "":
		baseResp.Data["name"] = name
		baseResp.Data["sig"] = sig
		baseResp.Data["expired"] = expired
	case appSig != "" && callback != "":
		SaveSessionAttrs(user.Id, appSig, int(time.Now().Unix())+3600*24*30)
		baseResp.Data["name"] = name
		baseResp.Data["sig"] = appSig
		baseResp.Data["expired"] = int(time.Now().Unix()) + 3600*24*30
	default:
		sig, expired := this.CreateSession(user.Id, 3600*24*30)
		baseResp.Data["name"] = name
		baseResp.Data["sig"] = sig
		baseResp.Data["expired"] = expired
	}
	this.ServeApiJson(baseResp)
	return
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
	repeatPassword := strings.TrimSpace(this.GetString("repeatPassword", ""))

	var lastID int64
	var err error
	switch {
	case password != repeatPassword:
		this.ResposeError(baseResp, "password not equal the repeart one")
		return
	case !utils.IsUsernameValid(name):
		this.ResposeError(baseResp, "name pattern is invalid")
		return
	case ReadUserIdByName(name) > 0:
		this.ResposeError(baseResp, "name is already existent")
		return
	default:
		lastID, err = InsertRegisterUser(name, str.Md5Encode(g.Config().Salt+password), email)
		if err != nil {
			this.ResposeError(baseResp, "insert user fail "+err.Error())
			return
		}
	}

	sig, expired := this.CreateSession(lastID, 3600*24*30)
	baseResp.Data["name"] = name
	baseResp.Data["sig"] = sig
	baseResp.Data["expired"] = expired
	this.ServeApiJson(baseResp)
	return
}

func (this *AuthApiController) GetUser() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()
	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	} else {
		username := this.GetString("cName", this.Ctx.GetCookie("name"))
		user := ReadUserByName(username)
		if user == nil {
			this.ResposeError(baseResp, "not found user")
			return
		}
		baseResp.Data["name"] = user.Name
		baseResp.Data["email"] = user.Email
		baseResp.Data["cnname"] = user.Cnname
		baseResp.Data["im"] = user.IM
		baseResp.Data["qq"] = user.QQ
		baseResp.Data["phone"] = user.Phone
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *AuthApiController) UpdateUser() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()

	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	} else {
		username := this.GetString("cName", this.Ctx.GetCookie("name"))
		user := ReadUserByName(username)
		user.Email = strings.TrimSpace(this.GetString("email", user.Email))
		user.Cnname = strings.TrimSpace(this.GetString("cnname", user.Cnname))
		user.IM = strings.TrimSpace(this.GetString("im", user.IM))
		user.QQ = strings.TrimSpace(this.GetString("qq", user.QQ))
		user.Phone = strings.TrimSpace(this.GetString("phone", user.Phone))
		passwdtmp := strings.TrimSpace(this.GetString("password", ""))
		oldpasswdtmp := strings.TrimSpace(this.GetString("oldpassword", ""))
		if passwdtmp != "" {
			if user.Passwd != str.Md5Encode(g.Config().Salt+oldpasswdtmp) {
				this.ResposeError(baseResp, "original password is empty or the password you inputed is not matched the original one.")
				return
			} else {
				user.Passwd = str.Md5Encode(g.Config().Salt + passwdtmp)
			}
		}
		_, err := user.Update()
		if err != nil {
			this.ResposeError(baseResp, err.Error())
			return
		}
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *AuthApiController) CountNumOfTeam() {
	baseResp := this.BasicRespGen()
	_, err := this.SessionCheck()

	if err != nil {
		this.ResposeError(baseResp, err.Error())
		return
	} else {
		numberOfteam, err := CountNumOfTeam()
		if err != nil {
			this.ResposeError(baseResp, err.Error())
			return
		}
		baseResp.Data["count"] = numberOfteam
	}
	this.ServeApiJson(baseResp)
	return
}

func (this *AuthApiController) CreateSession(uid int64, maxAge int) (sig string, expired int) {
	sig = utils.GenerateUUID()
	expired = int(time.Now().Unix()) + maxAge
	SaveSessionAttrs(uid, sig, expired)
	return
}
