package helper

import (
	"errors"

	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/uic"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	"github.com/gin-gonic/gin"
)

type WebSession struct {
	Name string
	Sig  string
}

func GetSession(c *gin.Context) (session WebSession, err error) {
	var name, sig string
	apiToken := c.Request.Header.Get("Apitoken")
	if apiToken == "" {
		err = errors.New("token key is not set")
		return
	}
	log.Debugf("header: %v, apiToken: %v", c.Request.Header, apiToken)
	var websession WebSession
	err = json.Unmarshal([]byte(apiToken), &websession)
	if err != nil {
		return
	}
	name = websession.Name
	log.Debugf("session got name: %s", name)
	if name == "" {
		err = errors.New("token key:name is empty")
		return
	}
	sig = websession.Sig
	log.Debugf("session got sig: %s", sig)
	if sig == "" {
		err = errors.New("token key:sig is empty")
		return
	}
	if err != nil {
		return
	}
	session = WebSession{name, sig}
	return
}

func SessionChecking(c *gin.Context) (auth bool, isServiceToken bool, err error) {
	auth = false
	var websessio WebSession
	websessio, err = GetSession(c)
	if err != nil {
		return
	}

	Serieves := config.ApiClient
	if Serieves.Enable && Serieves.NameIncludes(websessio.Name) {
		if Serieves.AuthToken(websessio.Name, websessio.Sig) {
			auth = true
			isServiceToken = true
			return
		}
		log.Warnf("use %s but got wrong sig (%s). Please need check this session", websessio.Name, websessio.Sig)
	}
	db := config.Con().Uic
	var user uic.User
	db.Where("name = ?", websessio.Name).Find(&user)
	if user.ID == 0 {
		err = errors.New("not found this user")
		return
	}
	var session uic.Session
	db.Table("session").Where("sig = ? and uid = ?", websessio.Sig, user.ID).Scan(&session)
	if session.ID == 0 {
		err = errors.New("session not found")
		return
	} else {
		auth = true
	}
	return
}

func GetUser(c *gin.Context) (user uic.User, err error) {
	db := config.Con().Uic
	websession, getserr := GetSession(c)
	if getserr != nil {
		err = getserr
		return
	}
	if v, ok := c.Get("is_service_token"); ok && v.(bool) {
		err = errors.New("services token no support this kind of action.")
		return
	}
	user = uic.User{
		Name: websession.Name,
	}
	dt := db.Where(&user).Find(&user)
	err = dt.Error
	return
}
