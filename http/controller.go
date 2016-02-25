package http

import (
	"fmt"
	"github.com/Cepave/alarm/g"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/toolkits/file"
	"log"

	// "sort"
	"strings"
	"time"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Version() {
	this.Ctx.WriteString(g.VERSION)
}

func (this *MainController) Health() {
	this.Ctx.WriteString("ok")
}

func (this *MainController) Workdir() {
	this.Ctx.WriteString(fmt.Sprintf("%s", file.SelfDir()))
}

func (this *MainController) ConfigReload() {
	remoteAddr := this.Ctx.Input.Context.Request.RemoteAddr
	if strings.HasPrefix(remoteAddr, "127.0.0.1") {
		g.ParseConfig(g.ConfigFile)
		this.Data["json"] = g.Config()
		this.ServeJSON()
	} else {
		this.Ctx.WriteString("no privilege")
	}
}

func SelectSessionBySig(sig string) *Session {
	if sig == "" {
		return nil
	}

	obj := Session{Sig: sig}
	err := orm.NewOrm().Read(&obj, "Sig")
	if err != nil {
		if err != orm.ErrNoRows {
			log.Println(err.Error())
		}
		return nil
	}
	return &obj
}

func DeleteSessionById(id int64) (int64, error) {
	r, err := orm.NewOrm().Raw("DELETE FROM `session` WHERE `id` = ?", id).Exec()
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

func SelectUserById(id int64) *User {
	if id <= 0 {
		return nil
	}

	obj := User{Id: id}
	err := orm.NewOrm().Read(&obj, "Id")
	if err != nil {
		if err != orm.ErrNoRows {
			log.Println(err.Error())
		}
		return nil
	}
	return &obj
}

/**
 * @function name:	func CheckLoginStatusByCookie(sig) bool
 * @description:	This function checks user's login status by value of "sig" cookie.
 * @related issues:	OWL-127
 * @param:			sig string
 * @return:			bool
 * @author:			Don Hsieh
 * @since:			10/15/2015
 * @last modified: 	10/30/2015
 * @called by:		func (this *MainController) Index()
 *					 in http/controller.go
 */
func CheckLoginStatusByCookie(sig string) bool {
	if sig == "" {
		return false
	}

	sessionObj := SelectSessionBySig(sig)
	if sessionObj == nil {
		log.Println("no such sig")
		return false
	}

	if int64(sessionObj.Expired) < time.Now().Unix() {
		log.Println("session expired")
		DeleteSessionById(sessionObj.Id)
		return false
	}

	user := SelectUserById(sessionObj.Uid)
	if user == nil {
		log.Println("no such user")
		return false
	}

	if len(user.Name) > 0 {
		return true
	} else {
		return false
	}
}

func (this *MainController) Index() {
	// Only check the Login status in Prodution env.
    if g.Config().Debug == false {
        checkLogin(this)
    }
	defer func() {
		this.Data["Now"] = time.Now().Unix()
		this.TplName = "index.html"
		this.Data["FalconPortal"] = g.Config().Shortcut.FalconPortal
		this.Data["FalconDashboard"] = g.Config().Shortcut.FalconDashboard
		this.Data["GrafanaDashboard"] = g.Config().Shortcut.GrafanaDashboard
		this.Data["FalconAlarm"] = g.Config().Shortcut.FalconAlarm
		this.Data["FalconUIC"] = g.Config().Shortcut.FalconUIC
	}()

	this.Data["Events"] = g.Events.CloneToOrderedEvents()
}

func (this *MainController) Event() {
    // Only check the Login status in Prodution env.
    if g.Config().Debug == false {
        checkLogin(this)
    }
    
    this.Data["json"] = g.Events.CloneToOrderedEvents()
    this.ServeJSON()
}

func (this *MainController) Solve() {
	ids := this.GetString("ids")
	if ids == "" {
		this.Ctx.WriteString("")
		return
	}

	idArr := strings.Split(ids, ",,")
	for i := 0; i < len(idArr); i++ {
		g.Events.Delete(idArr[i])
	}

	this.Ctx.WriteString("")
}

func checkLogin(m *MainController) {
    sig := m.Ctx.GetCookie("sig")
    isLoggedIn := CheckLoginStatusByCookie(sig)
    if !isLoggedIn {
        RedirectUrl := g.Config().RedirectUrl
        m.Redirect(RedirectUrl, 302)
    }
}

// func getOrderedEventsClone() []*g.EventDto {
// 	events := g.Events.Clone()
// 	count := len(events)
// 	if count == 0 {
// 		return []*g.EventDto{}
// 	}

// 	sortedEvent := make([]*g.EventDto, count)
// 	i := 0
// 	for _, event := range events {
// 		sortedEvent[i] = event
// 		i++
// 	}
//     // Sorted by Timestamp of EventDto
// 	sort.Sort(g.OrderedEvents(sortedEvent))
// 	return sortedEvent
// }
