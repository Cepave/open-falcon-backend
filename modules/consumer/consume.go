package main

import (
	"fmt"

	"github.com/Cepave/open-falcon-backend/modules/consumer/g"
	"github.com/Cepave/open-falcon-backend/modules/consumer/influx"
	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

func logOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func consume() {
	addr := fmt.Sprintf("amqp://%s:%s@%s/", g.Config().Mq.User, g.Config().Mq.Pass, g.Config().Mq.Addr)
	conn, err := amqp.Dial(addr)
	logOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	logOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"text", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	logOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	logOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Debugf("Received a message: %s", d.Body)
			influx.Send(d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
