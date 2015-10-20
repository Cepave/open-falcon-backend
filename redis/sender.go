package redis

import (
	"encoding/json"
	"github.com/Cepave/alarm/g"
	"github.com/Cepave/sender/model"
	"log"
	"strings"
)

func LPUSH(queue, message string) {
	rc := g.RedisConnPool.Get()
	defer rc.Close()
	_, err := rc.Do("LPUSH", queue, message)
	if err != nil {
		log.Println("LPUSH redis", queue, "fail:", err, "message:", message)
	}
}

func WriteSmsModel(sms *model.Sms) {
	if sms == nil {
		return
	}

	bs, err := json.Marshal(sms)
	if err != nil {
		log.Println(err)
		return
	}

	LPUSH(g.Config().Queue.Sms, string(bs))
}

func WriteMailModel(mail *model.Mail) {
	if mail == nil {
		return
	}

	bs, err := json.Marshal(mail)
	if err != nil {
		log.Println(err)
		return
	}

	LPUSH(g.Config().Queue.Mail, string(bs))
}

func WriteQQModel(qq *model.QQ) {
	if qq == nil {
		return
	}

	bs, err := json.Marshal(qq)
	if err != nil {
		log.Println(err)
		return
	}

	LPUSH(g.Config().Queue.QQ, string(bs))
}

func WriteSms(tos []string, content string) {
	if len(tos) == 0 {
		return
	}

	sms := &model.Sms{Tos: strings.Join(tos, ","), Content: content}
	WriteSmsModel(sms)
}

func WriteMail(tos []string, subject, content string) {
	if len(tos) == 0 {
		return
	}

	mail := &model.Mail{Tos: strings.Join(tos, ","), Subject: subject, Content: content}
	WriteMailModel(mail)
}

func WriteQQ(tos []string, subject, content string) {
	if len(tos) == 0 {
		return
	}

	qq := &model.QQ{Tos: strings.Join(tos, ","), Subject: subject, Content: content}
	WriteQQModel(qq)
}
