package cron

import (
	"encoding/json"
	"fmt"

	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/alarm/api"
	"github.com/Cepave/open-falcon-backend/modules/alarm/g"
	"github.com/Cepave/open-falcon-backend/modules/alarm/redis"
	"github.com/go-resty/resty"
	log "github.com/sirupsen/logrus"
)

func consume(event *model.Event, isHigh bool) {
	log.Info("consume")
	actionId := event.ActionId()
	if actionId <= 0 {
		return
	}

	forDemoSendEmail(event)
	action := api.GetAction(actionId)
	if action == nil {
		return
	}

	if action.Callback == 1 {
		Callback(event, action)
	}

	// if isHigh {
	// 	consumeHighEvents(event, action)
	// } else {
	// 	consumeLowEvents(event, action)
	// }
}

type PostMailDemo struct {
	Subject string `json:"subject"`
	Content string `json:"content"`
	TplId   int    `json:"tpl_id"`
}

func forDemoSendEmail(event *model.Event) {
	log.Info("forDemoSendEmail")
	f2econf := g.Config().F2eApiEmailHelper
	smsContent := GenerateSmsContent(event)
	mailContent := GenerateMailContent(event)
	postTmp := PostMailDemo{
		TplId:   event.TplId(),
		Subject: smsContent,
		Content: mailContent,
	}
	postbodyb, _ := json.Marshal(postTmp)
	Apitoken := fmt.Sprintf(`{"name": "%s", "sig": "%s"}`, f2econf.TokenName, f2econf.TokenKey)
	rt := resty.New()
	rt.SetHeader("Apitoken", Apitoken)
	resp, err := rt.R().
		SetBody(postTmp).
		Post(f2econf.URL)
	if err != nil {
		log.Errorf("send mail got error with: %v", err.Error())
	}
	log.Infof("send email got response: %v, postbody: %v", resp.String(), string(postbodyb))
}

// 高优先级的不做报警合并
func consumeHighEvents(event *model.Event, action *api.Action) {
	if action.Uic == "" {
		return
	}

	phones, mails := api.ParseTeams(action.Uic)

	smsContent := GenerateSmsContent(event)
	mailContent := GenerateMailContent(event)
	QQContent := GenerateQQContent(event)

	if event.Priority() < 3 {
		redis.WriteSms(phones, smsContent)
	}

	redis.WriteMail(mails, smsContent, mailContent)
	redis.WriteQQ(mails, smsContent, QQContent)
	ParseUserServerchan(event, action)
}

// 低优先级的做报警合并
func consumeLowEvents(event *model.Event, action *api.Action) {
	if action.Uic == "" {
		return
	}

	if event.Priority() < 3 {
		ParseUserSms(event, action)
	}

	ParseUserMail(event, action)
	ParseUserQQ(event, action)
	ParseUserServerchan(event, action)
}

func ParseUserSms(event *model.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)

	content := GenerateSmsContent(event)
	metric := event.Metric()
	status := event.Status
	priority := event.Priority()

	queue := g.Config().Redis.UserSmsQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := SmsDto{
			Priority: priority,
			Metric:   metric,
			Content:  content,
			Phone:    user.Phone,
			Status:   status,
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Println("json marshal SmsDto fail:", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Println("LPUSH redis", queue, "fail:", err, "dto:", string(bs))
		}
	}
}

func ParseUserMail(event *model.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)

	metric := event.Metric()
	subject := GenerateSmsContent(event)
	content := GenerateMailContent(event)
	status := event.Status
	priority := event.Priority()

	queue := g.Config().Redis.UserMailQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := MailDto{
			Priority: priority,
			Metric:   metric,
			Subject:  subject,
			Content:  content,
			Email:    user.Email,
			Status:   status,
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Println("json marshal MailDto fail:", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Println("LPUSH redis", queue, "fail:", err, "dto:", string(bs))
		}
	}
}

func ParseUserQQ(event *model.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)

	metric := event.Metric()
	subject := GenerateSmsContent(event)
	content := GenerateQQContent(event)
	status := event.Status
	priority := event.Priority()

	queue := g.Config().Redis.UserQQQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := QQDto{
			Priority: priority,
			Metric:   metric,
			Subject:  subject,
			Content:  content,
			Email:    user.Email,
			Status:   status,
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Println("json marshal QQDto fail:", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Println("LPUSH redis", queue, "fail:", err, "dto:", string(bs))
		}
	}
}

func ParseUserServerchan(event *model.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)
	metric := event.Metric()
	subject := GenerateSmsContent(event)
	content := GenerateServerchanContent(event)
	status := event.Status
	priority := event.Priority()

	queue := g.Config().Redis.UserServerchanQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := ServerchanDto{
			Priority: priority,
			Metric:   metric,
			Subject:  subject,
			Content:  content,
			Username: user.Name,
			Sckey:    user.IM,
			Status:   status,
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Println("json marshal ServerchanDto fail:", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Println("LPUSH redis", queue, "fail:", err, "dto:", string(bs))
		}
	}
}
