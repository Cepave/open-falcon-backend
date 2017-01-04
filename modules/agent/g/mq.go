package g

import (
	"encoding/json"
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/model"
	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

func logOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func SendToMQ(metrics []*model.MetricValue) {
	if len(metrics) == 0 {
		return
	}
	addr := fmt.Sprintf("amqp://%s:%s@%s/", Config().Mq.User, Config().Mq.Pass, Config().Mq.Addr)
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

	for _, val := range metrics {
		body, _ := json.Marshal(val)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		log.Debugln(" [x] Sent %s", body)
		logOnError(err, "Failed to publish a message")
	}
}
