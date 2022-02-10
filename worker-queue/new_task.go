package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:15672")
	if err != nil {
		panic(errors.Wrap(err, "failed to connect to rabbitmq server"))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(errors.Wrap(err, "failed to get channel"))
	}
	defer ch.Close()
	q, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)

	if err != nil {
		panic(errors.Wrap(err, "failed to declare queue"))
	}
	err = ch.Publish("", q.Name, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(os.Args[1]),
	})

}
