package mq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Cepave/fe/g"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
)

type Message struct {
	Hostname string `json:"hostname"`
	Mute     bool   `json:"mute"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func Start() {
	mq := g.Config().Mq
	nodes := map[string]interface{}{}
	conn, err := amqp.Dial(mq.Queue)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"mute", // name
		false,  // durable
		false,  // delete when usused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

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
