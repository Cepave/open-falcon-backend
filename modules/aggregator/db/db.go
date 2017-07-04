package db

import (
	"database/sql"
	"github.com/Cepave/open-falcon-backend/modules/aggregator/g"
	log "github.com/sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", g.Config().Database.Addr)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}

	DB.SetMaxIdleConns(g.Config().Database.Idle)

	err = DB.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
}
