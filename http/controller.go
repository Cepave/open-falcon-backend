package http

import (
	"log"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"

	"fmt"
	"github.com/Cepave/alarm/g"
	"github.com/astaxie/beego"
	"github.com/toolkits/file"
	"sort"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/open-falcon/alarm/g"
	"github.com/toolkits/file"
)

type Session struct {
	Id   int
	Sig string
	Expired string
}

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
	remoteAddr := this.Ctx.Request.RemoteAddr
	if strings.HasPrefix(remoteAddr, "127.0.0.1") {
		g.ParseConfig(g.ConfigFile)
		this.Data["json"] = g.Config()
		this.ServeJSON()
	} else {
		this.Ctx.WriteString("no privilege")
	}
}

/**
 * @function name:	func CheckLoginStatusByCookie(sig) bool
 * @description:	This function checks user's login status by value of "sig" cookie.
 * @related issues:	OWL-127
 * @param:			sig string
 * @return:			bool
 * @author:			Don Hsieh
 * @since:			10/15/2015
 * @last modified: 	10/16/2015
 * @called by:		func (this *MainController) Index()
 *					 in http/controller.go
 */
func CheckLoginStatusByCookie(sig string) bool {
	if sig == "" {
		return false
	}
	database := g.Config().Database.Db
	table := g.Config().Database.Table
	o := orm.NewOrm()
	o.Using(database)

	var session Session
	err := o.QueryTable(table).Filter("sig", sig).One(&session)
	if err == orm.ErrMultiRows {
		// Have multiple records
		log.Printf("Returned Multi Rows Not One")
		return false
	}
	if err == orm.ErrNoRows {
		// No result
		log.Printf("Not row found")
		return false
	}
	expiredTimeString := session.Expired
	expiredTimeInt, err := strconv.ParseInt(expiredTimeString, 10, 64)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	now := time.Now().Unix()
	expired := now > expiredTimeInt
	if !expired {
		return true
	} else {
		return false
	}
}

func (this *MainController) Index() {
	sig := this.Ctx.GetCookie("sig")
	log.Println("sig =", sig)
	isLoggedIn := CheckLoginStatusByCookie(sig)
	if !isLoggedIn {
		RedirectUrl := g.Config().RedirectUrl
		this.Redirect(RedirectUrl, 302)
	}

	events := g.Events.Clone()

	defer func() {
		this.Data["Now"] = time.Now().Unix()
		this.TplName = "index.html"
	}()

	count := len(events)
	if count == 0 {
		this.Data["Events"] = []*g.EventDto{}
		return
	}

	// 按照持续时间排序
	beforeOrder := make([]*g.EventDto, count)
	i := 0
	for _, event := range events {
		beforeOrder[i] = event
		i++
	}

	sort.Sort(g.OrderedEvents(beforeOrder))
	this.Data["Events"] = beforeOrder
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