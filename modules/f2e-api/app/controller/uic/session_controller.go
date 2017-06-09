package uic

import (
	"net/http"

	h "github.com/Cepave/open-falcon-backend/modules/f2e-api/app/helper"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/uic"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type APILoginInput struct {
	Name     string `json:"name"  form:"name" binding:"required"`
	Password string `json:"password"  form:"password" binding:"required"`
}

func Login(c *gin.Context) {
	inputs := APILoginInput{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, "name or password is blank")
		return
	}
	name := inputs.Name
	password := inputs.Password

	user := uic.User{
		Name: name,
	}
	db.Uic.Where(&user).Find(&user)
	switch {
	case user.Name == "":
		h.JSONR(c, badstatus, "no such user")
		return
	case user.Passwd != utils.HashIt(password):
		h.JSONR(c, badstatus, "password error")
		return
	}
	session := uic.Session{
		Uid: user.ID,
	}
	session = session.FindVaildSession()
	resp := struct {
		Sig   string `json:"sig,omitempty"`
		Name  string `json:"name,omitempty"`
		Admin bool   `json:"admin"`
	}{session.Sig, user.Name, user.IsAdmin()}
	h.JSONR(c, resp)
	return
}

func Logout(c *gin.Context) {
	wsession, err := h.GetSession(c)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	var session uic.Session
	var user uic.User
	db.Uic.Table("user").Where(uic.User{Name: wsession.Name}).Scan(&user)
	db.Uic.Table("session").Where("sig = ? AND uid = ?", wsession.Sig, user.ID).Scan(&session)

	if session.ID == 0 {
		h.JSONR(c, badstatus, "not found this kind of session in database.")
		return
	} else {
		r := db.Uic.Table("session").Delete(&session)
		if r.Error != nil {
			h.JSONR(c, badstatus, r.Error)
		}
		h.JSONR(c, "logout successful")
	}
	return
}

func AuthSession(c *gin.Context) {
	auth, isService, err := h.SessionChecking(c)
	if err != nil || auth != true {
		h.JSONR(c, http.StatusUnauthorized, err)
		return
	}
	if isService {
		h.JSONR(c, "session is vaild, it's servies token")
		return
	} else {
		h.JSONR(c, "session is vaild!")
		return
	}
}

func CreateRoot(c *gin.Context) {
	password := c.DefaultQuery("password", "")
	signupDisable := viper.GetBool("signup_disable")
	switch {
	case password == "":
		h.JSONR(c, badstatus, "password is empty, please check it")
		return
	case signupDisable:
		h.JSONR(c, badstatus, "sign up is not enabled, please contact administrator")
		return
	}
	password = utils.HashIt(password)
	user := uic.User{
		Name:   "root",
		Passwd: password,
	}
	dt := db.Uic.Table("user").Save(&user)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, "root created!")
	return
}
