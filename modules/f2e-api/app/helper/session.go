package helper

import (
	"errors"

	"encoding/json"

	log "github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/uic"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
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

func SessionChecking(c *gin.Context) (auth bool, err error) {
	auth = false
	var websessio WebSession
	websessio, err = GetSession(c)
	if err != nil {
		return
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
	websession, _ := GetSession(c)
	user = uic.User{
		Name: websession.Name,
	}
	dt := db.Where(&user).Find(&user)
	err = dt.Error
	return
}
