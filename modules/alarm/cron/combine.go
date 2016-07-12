package cron

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/Cepave/alarm/api"
	"github.com/Cepave/alarm/g"
	redi "github.com/Cepave/alarm/redis"
	"log"
	"strings"
	"time"
)

func CombineSms() {
	for {
		// 每分钟读取处理一次
		time.Sleep(time.Minute)
		combineSms()
	}
}

func CombineMail() {
	for {
		// 每分钟读取处理一次
		time.Sleep(time.Minute)
		combineMail()
	}
}

func CombineQQ() {
	for {
		time.Sleep(time.Minute)
		combineQQ()
	}
}

func CombineServerchan() {
	for {
		time.Sleep(time.Minute)
		combineServerchan()
	}
}

func combineMail() {
	dtos := popAllMailDto()
	count := len(dtos)
	if count == 0 {
		return
	}

	dtoMap := make(map[string][]*MailDto)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%d%s%s%s", dtos[i].Priority, dtos[i].Status, dtos[i].Email, dtos[i].Metric)
		if _, ok := dtoMap[key]; ok {
			dtoMap[key] = append(dtoMap[key], dtos[i])
		} else {
			dtoMap[key] = []*MailDto{dtos[i]}
		}
	}

	// 不要在这处理，继续写回redis，否则重启alarm很容易丢数据
	for _, arr := range dtoMap {
		size := len(arr)
		if size == 1 {
			redi.WriteMail([]string{arr[0].Email}, arr[0].Subject, arr[0].Content)
			continue
		}

		subject := fmt.Sprintf("[P%d][%s] %d %s", arr[0].Priority, arr[0].Status, size, arr[0].Metric)
		contentArr := make([]string, size)
		for i := 0; i < size; i++ {
			contentArr[i] = arr[i].Content
		}
		content := strings.Join(contentArr, "\r\n")

		redi.WriteMail([]string{arr[0].Email}, subject, content)
	}
}

func combineSms() {
	dtos := popAllSmsDto()
	count := len(dtos)
	if count == 0 {
		return
	}

	dtoMap := make(map[string][]*SmsDto)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%d%s%s%s", dtos[i].Priority, dtos[i].Status, dtos[i].Phone, dtos[i].Metric)
		if _, ok := dtoMap[key]; ok {
			dtoMap[key] = append(dtoMap[key], dtos[i])
		} else {
			dtoMap[key] = []*SmsDto{dtos[i]}
		}
	}

	for _, arr := range dtoMap {
		size := len(arr)
		if size == 1 {
			redi.WriteSms([]string{arr[0].Phone}, arr[0].Content)
			continue
		}

		// 把多个sms内容写入数据库，只给用户提供一个链接
		contentArr := make([]string, size)
		for i := 0; i < size; i++ {
			contentArr[i] = arr[i].Content
		}
		content := strings.Join(contentArr, ",,")

		first := arr[0].Content
		t := strings.Split(first, "][")
		eg := ""
		if len(t) >= 3 {
			eg = t[2]
		}

		path, err := api.LinkToSMS(content)
		sms := ""
		if err != nil || path == "" {
			sms = fmt.Sprintf("[P%d][%s] %d %s.  e.g. %s detail in email", arr[0].Priority, arr[0].Status, size, arr[0].Metric, eg)
			log.Println("get link fail", err)
		} else {
			links := g.Config().Api.Links
			sms = fmt.Sprintf("[P%d][%s] %d %s e.g. %s %s/%s ", arr[0].Priority, arr[0].Status, size, arr[0].Metric, eg, links, path)
		}

		redi.WriteSms([]string{arr[0].Phone}, sms)
	}
}

func combineQQ() {
	dtos := popAllQQDto()
	count := len(dtos)
	if count == 0 {
		return
	}

	dtoMap := make(map[string][]*QQDto)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%d%s%s%s", dtos[i].Priority, dtos[i].Status, dtos[i].Email, dtos[i].Metric)
		if _, ok := dtoMap[key]; ok {
			dtoMap[key] = append(dtoMap[key], dtos[i])
		} else {
			dtoMap[key] = []*QQDto{dtos[i]}
		}
	}

	// 不要在这处理，继续写回redis，否则重启alarm很容易丢数据
	for _, arr := range dtoMap {
		size := len(arr)
		if size == 1 {
			redi.WriteQQ([]string{arr[0].Email}, arr[0].Subject, arr[0].Content)
			continue
		}
		subject := fmt.Sprintf("[P%d][%s] %d %s", arr[0].Priority, arr[0].Status, size, arr[0].Metric)
		contentArr := make([]string, size)
		for i := 0; i < size; i++ {
			contentArr[i] = arr[i].Content
		}
		content := strings.Join(contentArr, "\r\n")
		redi.WriteQQ([]string{arr[0].Email}, subject, content)
	}
}

func combineServerchan() {
	dtos := popAllServerchanDto()
	count := len(dtos)
	if count == 0 {
		return
	}

	dtoMap := make(map[string][]*ServerchanDto)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%d%s%s%s", dtos[i].Priority, dtos[i].Status, dtos[i].Sckey, dtos[i].Metric)
		if _, ok := dtoMap[key]; ok {
			dtoMap[key] = append(dtoMap[key], dtos[i])
		} else {
			dtoMap[key] = []*ServerchanDto{dtos[i]}
		}
	}

	for _, arr := range dtoMap {
		size := len(arr)
		if size == 1 {
			redi.WriteServerchan([]string{arr[0].Sckey}, arr[0].Subject, arr[0].Content)
			continue
		}
		subject := fmt.Sprintf("[P%d][%s] %d %s", arr[0].Priority, arr[0].Status, size, arr[0].Metric)
		contentArr := make([]string, size)
		for i := 0; i < size; i++ {
			contentArr[i] = arr[i].Content
		}
		content := strings.Join(contentArr, "\r\n")
		redi.WriteServerchan([]string{arr[0].Sckey}, subject, content)
	}
}

func popAllSmsDto() []*SmsDto {
	ret := []*SmsDto{}
	queue := g.Config().Redis.UserSmsQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Println("get SmsDto fail", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var smsDto SmsDto
		err = json.Unmarshal([]byte(reply), &smsDto)
		if err != nil {
			log.Printf("json unmarshal SmsDto: %s fail: %v", reply, err)
			continue
		}

		ret = append(ret, &smsDto)
	}

	return ret
}

func popAllMailDto() []*MailDto {
	ret := []*MailDto{}
	queue := g.Config().Redis.UserMailQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Println("get MailDto fail", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var mailDto MailDto
		err = json.Unmarshal([]byte(reply), &mailDto)
		if err != nil {
			log.Printf("json unmarshal MailDto: %s fail: %v", reply, err)
			continue
		}

		ret = append(ret, &mailDto)
	}

	return ret
}

func popAllQQDto() []*QQDto {
	ret := []*QQDto{}
	queue := g.Config().Redis.UserQQQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Println("get QQDto fail", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var qqDto QQDto
		err = json.Unmarshal([]byte(reply), &qqDto)
		if err != nil {
			log.Printf("json unmarshal QQDto: %s fail: %v", reply, err)
			continue
		}

		ret = append(ret, &qqDto)
	}

	return ret
}

func popAllServerchanDto() []*ServerchanDto {
	ret := []*ServerchanDto{}
	queue := g.Config().Redis.UserServerchanQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Println("get ServerchanDto fail", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var serverchanDto ServerchanDto
		err = json.Unmarshal([]byte(reply), &serverchanDto)
		if err != nil {
			log.Printf("json unmarshal ServerchanDto: %s fail: %v", reply, err)
			continue
		}

		ret = append(ret, &serverchanDto)
	}

	return ret
}
