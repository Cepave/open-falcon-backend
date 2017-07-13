package uic

import (
	"errors"
	"fmt"

	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/uic"
	"github.com/gin-gonic/gin"
	"github.com/masato25/resty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type BossUserInfoOutputStruct struct {
	Status int
	Info   string
	Data   struct {
		UserID       string        `json:"user_id"`
		UserName     string        `json:"username"`
		WeChat       string        `json:"wechat"`
		Cell         string        `json:"cell"`
		Email        string        `json:"email"`
		Telphone     string        `json:"telphone"`
		Realname     string        `json:"realname"`
		Position     string        `json:"position"`
		DepartmentID string        `json:"department_id"`
		TeamID       []interface{} `json:"team_id"`
		Roles        []string      `json:"roles"`
	}
}

func GetBossUserInfoByCookie(c *gin.Context) {
	bossToken, err := c.Cookie("FASTWEB_CDNBOSS_SEESION")
	if bossToken == "" {
		h.JSONR(c, badstatus, "FASTWEB_CDNBOSS_SEESION is invaild")
		return
	} else if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	bson, err := fetchBossUserInfo(bossToken)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	log.Debugf("debug boss get_user_info response:%v", bson)
	user := uic.User{
		Name: bson.Data.UserName,
	}
	userExist := user.UserNameExist()
	if !userExist {
		user, err = createNewUser(bson)
		if err != nil {
			h.JSONR(c, badstatus, err.Error())
			return
		}
	} else {
		if dt := db.Uic.Model(&user).Where("name = ?", user.Name).Scan(&user); dt.Error != nil {
			h.JSONR(c, badstatus, dt.Error)
			return
		}
		// user disalbed
		if user.Role == -1 {
			h.JSONR(c, badstatus, fmt.Errorf("user: %s has been disabled", user.Name))
			return
		}
	}
	usession := uic.Session{
		Uid: user.ID,
	}
	// find a vaild session and return
	usession = usession.FindVaildSession()
	h.JSONR(c, struct {
		Sig   string `json:"sig,omitempty"`
		Name  string `json:"name,omitempty"`
		Admin bool   `json:"admin"`
	}{usession.Sig, user.Name, user.IsAdmin()})
}

func fetchBossUserInfo(bossToken string) (bson BossUserInfoOutputStruct, err error) {
	bson = BossUserInfoOutputStruct{}
	rt := resty.New()
	rt.SetCookie(&http.Cookie{
		Name:  "FASTWEB_CDNBOSS_SEESION",
		Value: bossToken})
	bossAPIURI := viper.GetString("boss_api.boss_get_info_v2")
	if bossAPIURI == "" {
		err = errors.New("boss_api.boss_get_info_v2 is not set, please contact admin")
		return
	}
	resp, err0 := rt.R().Get(bossAPIURI)
	if err0 != nil {
		err = err0
		return
	}
	err = json.Unmarshal(resp.Body(), &bson)
	return
}

func createNewUser(bson BossUserInfoOutputStruct) (user uic.User, err error) {
	tx := db.Uic.Begin()
	if bson.Data.UserName == "" {
		err = errors.New("username can not be empty")
		return
	}
	user = uic.User{
		Name:   bson.Data.UserName,
		Passwd: "",
		Cnname: bson.Data.Realname,
		Email:  bson.Data.Email,
		Phone:  bson.Data.Cell,
		IM:     bson.Data.WeChat,
		Role:   0,
		//boss create set creator = 1
		Creator: 1,
	}
	if dt := tx.Model(&user).Save(&user); dt.Error != nil {
		tx.Rollback()
		err = dt.Error
		return
	}
	tx.Commit()
	return
}

// Deprecated
func ForwardToBossLoginPage(c *gin.Context) {
	referFrom := c.Request.Header.Get("Referer")
	if referFrom == "" {
		h.JSONR(c, badstatus, errors.New("no referer site info"))
		return
	}
	bossFToken := viper.GetString("boss_api.boss_f_token")
	bossHost := viper.GetString("boss_api.boss_host")
	if bossFToken == "" || bossHost == "" {
		h.JSONR(c, badstatus, errors.New("boss_f_token or boss_host is not set"))
		return
	}
	currentHost := fmt.Sprintf("%s/api/v1/third-party/auth/login/__TOKEN__?callback=%s", c.Request.Host, referFrom)
	hashedKey := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s||%s", bossFToken, currentHost)))
	forwardToUriEncode := url.QueryEscape(hashedKey)
	forwardToUri := fmt.Sprintf("%s/Third/login/%s", bossHost, forwardToUriEncode)
	h.JSONR(c, map[string]interface{}{
		"a": forwardToUri,
		"b": fmt.Sprintf("%s||%s", bossFToken, currentHost),
	})
}

// Deprecated
func BossRedirectLogin(c *gin.Context) {
	utoken := c.Params.ByName("utoken")
	callBack := c.Query("callback")
	if utoken == "" || callBack == "" {
		h.JSONR(c, badstatus, "bad request")
		return
	}
	queryKey := viper.GetString("boss_api.user_query_key")
	h.JSONR(c, map[string]interface{}{
		"callback": callBack,
		"utoken":   utoken,
		"qkey":     queryKey,
	})
}

// Deprecated
func RedirectToOriginalPage(c *gin.Context) {
	referFrom := c.Request.Header.Get("Referer")
	if referFrom == "" {
		h.JSONR(c, badstatus, errors.New("no referer info"))
	}
}
