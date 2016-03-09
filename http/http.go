package http

import (
	"time"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"github.com/Cepave/alarm/g"
	"github.com/astaxie/beego"
	"log"
	_ "net/http/pprof"

	"github.com/astaxie/beego"
	"github.com/open-falcon/alarm/g"
)

func configRoutes() {
	beego.Router("/", &MainController{}, "get:Index")
	beego.Router("/version", &MainController{}, "get:Version")
	beego.Router("/health", &MainController{}, "get:Health")
	beego.Router("/workdir", &MainController{}, "get:Workdir")
	beego.Router("/config/reload", &MainController{}, "get:ConfigReload")
    beego.Router("/event", &MainController{}, "get:Event")
	beego.Router("/event/solve", &MainController{}, "post:Solve")
}

func Duration(now, before int64) string {
	d := now - before
	if d <= 60 {
		return "just now"
	}

	if d <= 120 {
		return "1 minute ago"
	}

	if d <= 3600 {
		return fmt.Sprintf("%d minutes ago", d/60)
	}

	if d <= 7200 {
		return "1 hour ago"
	}
	
	if d <= 3600*24 {
		return fmt.Sprintf("%d hours ago", d/3600)
	}

	if d <= 3600*24*2 {
		return "1 day ago"
	}

	return fmt.Sprintf("%d days ago", d/3600/24)
}

func init() {
	configRoutes()
	beego.AddFuncMap("duration", Duration)
}

type User struct {
	Id      int64     `json:"id"`
	Name    string    `json:"name"`
	Cnname  string    `json:"cnname"`
	Passwd  string    `json:"-"`
	Email   string    `json:"email"`
	Phone   string    `json:"phone"`
	IM      string    `json:"im" orm:"column(im)"`
	QQ      string    `json:"qq" orm:"column(qq)"`
	Role    int       `json:"role"`
	Created time.Time `json:"-" orm:"-"`
}

type Session struct {
	Id      int64
	Uid     int64
	Sig     string
	Expired int
}

func InitDatabase() {
	// set default database
	config := g.Config()
	orm.RegisterDataBase("default", "mysql", config.Uic.Addr, config.Uic.Idle, config.Uic.Max)

	// register model
	orm.RegisterModel(new(User), new(Session))

	if config.Debug == true {
		orm.Debug = true
	}
}

func Start() {
	if !g.Config().Http.Enabled {
		return
	}

	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}

	if g.Config().Debug {
		beego.BConfig.RunMode = "dev"
	} else {
		beego.BConfig.RunMode = "prod"
	}

	InitDatabase()

	beego.Run(addr)

	log.Println("http listening", addr)
}
