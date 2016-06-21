package uic

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Cepave/fe/g"
	"github.com/Cepave/fe/http/base"
	. "github.com/Cepave/fe/model/uic"
	"github.com/Cepave/fe/utils"
	"github.com/toolkits/str"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AuthController struct {
	base.BaseController
}

func (this *AuthController) Logout() {
	u := this.Ctx.Input.GetData("CurrentUser").(*User)
	token := this.Ctx.GetCookie("token")
	if len(token) > 0 {
		url := g.Config().Api.Logout + "/" + token
		log.Println("logout url =", url)
		result := sendHttpGetRequest(url)
		log.Println("logout result =", result)
		this.Ctx.SetCookie("token", "", 0, "/")
		this.Ctx.SetCookie("token", "", 0, "/", g.Config().Http.Cookie)
	}
	RemoveSessionByUid(u.Id)
	this.Ctx.SetCookie("sig", "", 0, "/")
	this.Ctx.SetCookie("sig", "", 0, "/", g.Config().Http.Cookie)
	this.Ctx.SetCookie("name", "", 0, "/")
	this.Ctx.SetCookie("name", "", 0, "/", g.Config().Http.Cookie)
	this.Redirect("/auth/login", 302)
}

func (this *AuthController) LoginGet() {
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

func (this *AuthController) LoginPost() {
	name := this.GetString("name", "")
	password := this.GetString("password", "")
	if name == "" || password == "" {
		this.ServeErrJson("name or password is blank")
		return
	}

	var u *User

	ldapEnabled := this.MustGetBool("ldap", false)

	if ldapEnabled {
		sucess, err := utils.LdapBind(g.Config().Ldap.Addr,
			g.Config().Ldap.BaseDN,
			g.Config().Ldap.BindDN,
			g.Config().Ldap.BindPasswd,
			g.Config().Ldap.UserField,
			name,
			password)
		if err != nil {
			this.ServeErrJson(err.Error())
			return
		}

		if !sucess {
			this.ServeErrJson("name or password error")
			return
		}

		user_attributes, err := utils.Ldapsearch(g.Config().Ldap.Addr,
			g.Config().Ldap.BaseDN,
			g.Config().Ldap.BindDN,
			g.Config().Ldap.BindPasswd,
			g.Config().Ldap.UserField,
			name,
			g.Config().Ldap.Attributes)
		userSn := ""
		userMail := ""
		userTel := ""
		if err == nil {
			userSn = user_attributes["sn"]
			userMail = user_attributes["mail"]
			userTel = user_attributes["telephoneNumber"]
		}

		arr := strings.Split(name, "@")
		var userName, userEmail string
		if len(arr) == 2 {
			userName = arr[0]
			userEmail = name
		} else {
			userName = name
			userEmail = userMail
		}

		u = ReadUserByName(userName)
		if u == nil {
			// 说明用户不存在
			u = &User{
				Name:   userName,
				Passwd: "",
				Cnname: userSn,
				Phone:  userTel,
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

func (this *AuthController) renderLoginPage(sig, callback string) {
	this.Data["CanRegister"] = g.Config().CanRegister
	this.Data["LdapEnabled"] = g.Config().Ldap.Enabled
	this.Data["Sig"] = sig
	this.Data["Callback"] = callback
	this.Data["Shortcut"] = g.Config().Shortcut
	this.TplName = "auth/login.html"
}

func (this *AuthController) RegisterGet() {
	this.Data["CanRegister"] = g.Config().CanRegister
	this.Data["Shortcut"] = g.Config().Shortcut
	this.TplName = "auth/register.html"
}

func (this *AuthController) RegisterPost() {
	if !g.Config().CanRegister {
		this.ServeErrJson("registration system is not open")
		return
	}

	name := strings.TrimSpace(this.GetString("name", ""))
	password := strings.TrimSpace(this.GetString("password", ""))
	repeatPassword := strings.TrimSpace(this.GetString("repeat_password", ""))
	email := strings.TrimSpace(this.GetString("email", ""))

	if password != repeatPassword {
		this.ServeErrJson("password not equal the repeart one")
		return
	}

	if !utils.IsUsernameValid(name) {
		this.ServeErrJson("name pattern is invalid")
		return
	}

	if ReadUserIdByName(name) > 0 {
		this.ServeErrJson("name is already existent")
		return
	}

	lastId, err := InsertRegisterUser(name, str.Md5Encode(g.Config().Salt+password), email)
	if err != nil {
		this.ServeErrJson("insert user fail " + err.Error())
		return
	}

	this.CreateSession(lastId, 3600*24*30)

	this.ServeOKJson()
}

func (this *AuthController) CreateSession(uid int64, maxAge int) int {
	sig := utils.GenerateUUID()
	user := SelectUserById(uid)
	expired := int(time.Now().Unix()) + maxAge
	SaveSessionAttrs(uid, sig, expired)
	this.Ctx.SetCookie("sig", sig, maxAge, "/")
	this.Ctx.SetCookie("sig", sig, maxAge, "/", g.Config().Http.Cookie)
	this.Ctx.SetCookie("name", user.Name, maxAge, "/")
	this.Ctx.SetCookie("name", user.Name, maxAge, "/", g.Config().Http.Cookie)
	return expired
}

/**
 * @function name:   func (this *AuthController) LoginThirdParty()
 * @description:     This function returns third party login URL.
 * @related issues:  OWL-206
 * @param:           void
 * @return:          void
 * @author:          Don Hsieh
 * @since:           12/17/2015
 * @last modified:   12/17/2015
 * @called by:       beego.Router("/auth/third-party", &AuthController{}, "post:LoginThirdParty")
 *                    in fe/http/uic/uic_routes.go
 */
func (this *AuthController) LoginThirdParty() {
	s := g.Config().Api.Redirect
	s = base64.StdEncoding.EncodeToString([]byte(s))
	strEncoded := url.QueryEscape(s)
	loginUrl := g.Config().Api.Login + "/" + strEncoded
	this.ServeDataJson(loginUrl)
}

/**
 * @function name:   func sendHttpGetRequest(url string) map[string]interface{}
 * @description:     This function sends GET request to given URL.
 * @related issues:  OWL-206, OWL-159
 * @param:           url string
 * @return:          map[string]interface{}
 * @author:          Don Hsieh
 * @since:           12/17/2015
 * @last modified:   12/17/2015
 * @called by:       func (this *AuthController) LoginWithToken()
 *                    in fe/http/uic/auth_controller.go
 */
func sendHttpGetRequest(url string) map[string]interface{} {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error =", err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println("Error =", err.Error())
		return nil
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		var nodes = make(map[string]interface{})
		if err := json.Unmarshal(body, &nodes); err != nil {
			log.Println("Error =", err.Error())
			return nil
		}
		return nodes
	}
}

func setUserInfo(nodes map[string]interface{}, userInfo map[string]string) {
	if status, ok := nodes["status"]; ok {
		if int(status.(float64)) == 1 {
			data := nodes["data"].(map[string]interface{})
			access_key := data["access_key"].(string)
			username := data["username"].(string)
			email := data["email"].(string)
			log.Println("access_key =", access_key)
			userInfo["username"] = username
			userInfo["email"] = email
			userInfo["access_key"] = access_key
		}
	}
}

func getUserRole(access_key string) int {
	urlRole := g.Config().Api.Role + "/" + access_key
	nodes := sendHttpGetRequest(urlRole)
	role := -1
	if int(nodes["status"].(float64)) == 1 {
		permission := nodes["data"]
		log.Println("permission =", permission)
		if permission == "admin" {
			role = 0
		} else if permission == "operator" {
			role = 1
		} else if permission == "observer" {
			role = 2
		} else if permission == "deny" {
			role = 3
		}
	}

	// TODO: The role should be able to be changed on BOSS in the future.
	role = 1
	return role
}

/**
 * @function name:   func (this *AuthController) LoginWithToken()
 * @description:     This function logins user with third party token.
 * @related issues:  OWL-247, OWL-206
 * @param:           void
 * @return:          void
 * @author:          Don Hsieh
 * @since:           12/16/2015
 * @last modified:   01/08/2016
 * @called by:       beego.Router("/auth/login/:token", &AuthController{}, "get:LoginWithToken")
 *                    in fe/http/uic/uic_routes.go
 */
func (this *AuthController) LoginWithToken() {
	log.Println("func (this *AuthController) LoginWithToken()")
	token := this.Ctx.Input.Param(":token")
	log.Println("token =", token)
	key := g.Config().Api.Key
	authUrl := g.Config().Api.Access + "/" + token + "/" + key

	nodes := sendHttpGetRequest(authUrl)
	if nodes == nil {
		nodes = sendHttpGetRequest(authUrl)
	}
	log.Println("nodes =", nodes)

	var userInfo = make(map[string]string)
	userInfo["username"] = ""
	userInfo["email"] = ""
	userInfo["access_key"] = ""
	if nodes != nil {
		setUserInfo(nodes, userInfo)
	}
	log.Println("userInfo =", userInfo)

	username := userInfo["username"]
	if len(username) > 0 {
		access_key := userInfo["access_key"]
		user := ReadUserByName(username)
		if user == nil { // create third party user
			InsertRegisterUser(username, "", "")
			user = ReadUserByName(username)
		}
		if len(user.Passwd) == 0 {
			role := getUserRole(access_key)
			if role < 1 {
				role = getUserRole(access_key)
			}
			email := userInfo["email"]
			user.Email = email
			user.Role = role
			user.Update()
		}
		maxAge := 3600 * 24 * 30
		this.Ctx.SetCookie("token", access_key, maxAge, "/")
		this.Ctx.SetCookie("token", access_key, maxAge, "/", g.Config().Http.Cookie)

		appSig := this.GetString("sig", "")
		callback := this.GetString("callback", "")
		if appSig != "" && callback != "" {
			SaveSessionAttrs(user.Id, appSig, int(time.Now().Unix())+3600*24*30)
		} else {
			this.CreateSession(user.Id, 3600*24*30)
		}
		this.Redirect("/me/info", 302)
	} else {
		// not logged in. redirect to login page.
		appSig := this.GetString("sig", "")
		callback := this.GetString("callback", "")
		this.renderLoginPage(appSig, callback)
	}
}
