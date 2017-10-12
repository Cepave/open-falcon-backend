package hbsdb

import (
	"fmt"

	"github.com/Cepave/open-falcon-backend/common/model"
	log "github.com/sirupsen/logrus"
)

func QueryConfig(key string) (*model.Config, error) {
	sql := fmt.Sprintf("select t.key, t.value from common_config as t where t.key = '%s'", key)
	row := DB.QueryRow(sql)

	e := model.Config{}
	err := row.Scan(&e.Key, &e.Value)

	if err != nil {
		log.Warnln(err)
	}

	ret := &e
	return ret, nil
}
