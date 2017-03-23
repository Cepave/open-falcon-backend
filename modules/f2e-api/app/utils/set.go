package utils

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func UniqSet(list []interface{}) []interface{} {
	mapTmp := map[interface{}]int{}
	for _, s := range list {
		if _, ok := mapTmp[s]; !ok {
			mapTmp[s] = 1
		}
	}
	var keys []interface{}
	for k, _ := range mapTmp {
		keys = append(keys, k)
	}
	return keys
}

func MapTake(list []interface{}, limit int) []interface{} {
	res := []interface{}{}
	if limit > len(list) {
		limit = len(list)
	}
	for i := 0; i < limit; i++ {
		res = append(res, list[i])
	}
	return res
}

func ConverIntStringToList(eid string) (result string) {
	for i, e := range strings.Split(eid, ",") {
		v, err := strconv.Atoi(e)
		if err != nil {
			log.Debug(err.Error())
		} else {
			if i == 0 {
				result = fmt.Sprintf("%d", v)
			} else {
				result = fmt.Sprintf("%s, %d", result, v)
			}
		}
	}
	return
}
