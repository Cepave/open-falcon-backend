package mq

import (
	"bytes"
	"encoding/json"
	"github.com/Cepave/fe/g"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var RetriedLimit = 10
var SleepTimePeriod = time.Duration(60)
var ExitStringPrefix = "Exit mq goroutine safely: "
var LogStringFormat = "%s: %s"

type Message struct {
	Hostname string `json:"hostname"`
	Mute     bool   `json:"mute"`
}

func setup(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf(LogStringFormat, "Failed to connect to RabbitMQ", err)
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf(LogStringFormat, "Failed to open a channel", err)
		defer conn.Close()
		return nil, nil, err
	}
	return conn, ch, nil
}

func Start() {
	mq := g.Config().Mq
	nodes := map[string]interface{}{}

	// Retry RetriedLimit times if there is some problem during connecting
	var ch *amqp.Channel
	var conn *amqp.Connection
	var err error
	for i := 0; i < RetriedLimit; i++ {
		if conn, ch, err = setup(mq.Queue); err != nil {
			time.Sleep(time.Second * SleepTimePeriod)
		}
	}
	if err != nil {
		log.Println(ExitStringPrefix + "retried too many times.")
		return
	}
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"mute", // name
		false,  // durable
		false,  // delete when usused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		log.Printf(LogStringFormat, ExitStringPrefix+"failed to declare a queue", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Printf(LogStringFormat, ExitStringPrefix+"failed to register a consumer", err)
		return
	}

	forever := make(chan bool)

	var pkt Message
	go func() {
		for d := range msgs {
			json.Unmarshal(d.Body, &pkt)
			log.Printf("Received a message: %v", pkt)
			params := map[string]string{
				"host": pkt.Hostname,
				"mute": "",
			}
			if pkt.Mute {
				params["mute"] = "1"
			} else {
				params["mute"] = "0"
			}
			payload := map[string]interface{}{
				"method": "host.update",
				"params": params,
			}
			log.Println("payload =", payload)

			s, err := json.Marshal(payload)
			if err != nil {
				log.Println("json.Marshal Error:", err.Error())
			}

			url := mq.Consumer
			log.Println("url =", url)
			reqPost, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(s)))
			if err != nil {
				log.Println("reqPost Error:", err.Error())
			}
			reqPost.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(reqPost)
			if err != nil {
				log.Println("resp Error:", err.Error())
			} else {
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				err = json.Unmarshal(body, &nodes)
				if err != nil {
					log.Println("Unmarshal Error:", err.Error())
				}
				log.Println("nodes =", nodes)
			}
		}
	}()
	log.Println("mq.Start ok. Waiting for messages.")
	<-forever
}
