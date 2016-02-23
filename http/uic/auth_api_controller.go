package uic

import (
	"encoding/base64"
	"github.com/Cepave/fe/g"
	"github.com/Cepave/fe/http/base"
	. "github.com/Cepave/fe/model/uic"
	"github.com/Cepave/fe/utils"
	"github.com/toolkits/str"
	"log"
	"net/url"
	"strings"
	"time"
)

type AuthApiController struct {
	base.BaseController
}

func (this *AuthApiController) Logout() {
	u := this.Ctx.Input.GetData("CurrentUser").(*User)
	token := this.Ctx.GetCookie("token")
	if len(token) > 0 {
		url := g.Config().Api.Logout + "/" + token
		log.Println("logout url =", url)
		result := getRequest(url)
		log.Println("logout result =", result)
		this.Ctx.SetCookie("token", "", 0, "/")
		this.Ctx.SetCookie("token", "", 0, "/", g.Config().Http.Cookie)
	}
	RemoveSessionByUid(u.Id)
	this.Ctx.SetCookie("sig", "", 0, "/")
	this.Ctx.SetCookie("sig", "", 0, "/", g.Config().Http.Cookie)
	this.Redirect("/auth/login", 302)
}

func (this *AuthApiController) LoginGet() {
	appSig := this.GetString("sig", "")
	callback := this.GetString("callback", "")

	cookieSig := this.Ctx.GetCookie("sig")
	if cookieSig == "" {
		this.renderLoginPage(appSig, callback)
		return
	}

	sessionObj := ReadSessionBySig(cookieSig)
	if sessionObj == nil {
		this.renderLoginPage(appSig, callback)
		return
	}

	if int64(sessionObj.Expired) < time.Now().Unix() {
		RemoveSessionByUid(sessionObj.Uid)
		this.renderLoginPage(appSig, callback)
		return
	}

	if appSig != "" && callback != "" {
		this.Redirect(callback, 302)
	} else {
		this.Redirect("/me/info", 302)
	}
}

func (this *AuthApiController) LoginPost() {
	name := this.GetString("name", "")
	password := this.GetString("password", "")

	if name == "" || password == "" {
		this.ServeErrJson("name or password is blank")
		return
	}

	var u *User

	ldapEnabled := this.MustGetBool("ldap", false)

	if ldapEnabled {
		sucess, err := utils.LdapBind(g.Config().Ldap.Addr, name, password)
		if err != nil {
			this.ServeErrJson(err.Error())
			return
		}

		if !sucess {
			this.ServeErrJson("name or password error")
			return
		}

		arr := strings.Split(name, "@")
		var userName, userEmail string
		if len(arr) == 2 {
			userName = arr[0]
			userEmail = name
		} else {
			userName = name
			userEmail = ""
		}

		u = ReadUserByName(userName)
		if u == nil {
			// 说明用户不存在
			u = &User{
				Name:   userName,
				Passwd: "",
				Email:  userEmail,
			}
			_, err = u.Save()
			if err != nil {
				this.ServeErrJson("insert user fail " + err.Error())
				return
			}
		}
	} else {
		u = ReadUserByName(name)
		if u == nil {
			this.ServeErrJson("no such user")
			return
		}

		if u.Passwd != str.Md5Encode(g.Config().Salt+password) {
			this.ServeErrJson("password error")
			return
		}
	}

	appSig := this.GetString("sig", "")
	callback := this.GetString("callback", "")
	if appSig != "" && callback != "" {
		SaveSessionAttrs(u.Id, appSig, int(time.Now().Unix())+3600*24*30)
	} else {
		this.CreateSession(u.Id, 3600*24*30)
	}

	this.ServeDataJson(callback)
}

func (this *AuthApiController) renderLoginPage(sig, callback string) {
	this.Data["CanRegister"] = g.Config().CanRegister
	this.Data["LdapEnabled"] = g.Config().Ldap.Enabled
	this.Data["Sig"] = sig
	this.Data["Callback"] = callback
	this.Data["Shortcut"] = g.Config().Shortcut
	this.TplName = "auth/login.html"
}

func (this *AuthApiController) RegisterGet() {
	this.Data["CanRegister"] = g.Config().CanRegister
	this.Data["Shortcut"] = g.Config().Shortcut
	this.TplName = "auth/register.html"
}

func (this *AuthApiController) RegisterPost() {
	apiBasicParams, errorTmp, dataTmp := this.ApiBasicParams()
	apiBasicParams.Method = "POST"
	if !g.Config().CanRegister {
		errorTmp["message"] = "registration system is not open"
		apiBasicParams.Error = errorTmp
		this.ServeApiJson(apiBasicParams)
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
		errorTmp["message"] = "password not equal the repeart one"
	case !utils.IsUsernameValid(name):
		errorTmp["message"] = "name pattern is invalid"
	case ReadUserIdByName(name) > 0:
		errorTmp["message"] = "name is already existent"
	default:
		lastID, err = InsertRegisterUser(name, str.Md5Encode(g.Config().Salt+password), email)
		if err != nil {
			errorTmp["message"] = "insert user fail " + err.Error()
		}
	}
	if len(errorTmp) != 0 {
		apiBasicParams.Error = errorTmp
		this.ServeApiJson(apiBasicParams)
		return
	}

	sig, expired := this.CreateSession(lastID, 3600*24*30)
	dataTmp["sig"] = sig
	dataTmp["expired"] = expired
	apiBasicParams.Data = dataTmp
	this.ServeApiJson(apiBasicParams)
}

func (this *AuthApiController) CreateSession(uid int64, maxAge int) (sig string, expired int) {
	sig = utils.GenerateUUID()
	expired = int(time.Now().Unix()) + maxAge
	SaveSessionAttrs(uid, sig, expired)
	this.Ctx.SetCookie("sig", sig, maxAge, "/")
	this.Ctx.SetCookie("sig", sig, maxAge, "/", g.Config().Http.Cookie)
	return
}

/**
 * @function name:   func (this *AuthApiController) LoginThirdParty()
 * @description:     This function returns third party login URL.
 * @related issues:  OWL-206
 * @param:           void
 * @return:          void
 * @author:          Don Hsieh
 * @since:           12/17/2015
 * @last modified:   12/17/2015
 * @called by:       beego.Router("/auth/third-party", &AuthApiController{}, "post:LoginThirdParty")
 *                    in fe/http/uic/uic_routes.go
 */
func (this *AuthApiController) LoginThirdParty() {
	s := g.Config().Api.Redirect
	s = base64.StdEncoding.EncodeToString([]byte(s))
	strEncoded := url.QueryEscape(s)
	loginUrl := g.Config().Api.Login + "/" + strEncoded
	this.ServeDataJson(loginUrl)
}
