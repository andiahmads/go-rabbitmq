package main

import (
	"log"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	// conn, err := amqp091.Dial("amqp://guest:guest@localhost:15672")
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(errors.Wrap(err, "failed to connect to rabbitmq server"))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(errors.Wrap(err, "failed to get channel"))
	}
	defer ch.Close()
	q, err := ch.QueueDeclare("", false, false, true, false, nil)

	if err != nil {
		panic(errors.Wrap(err, "failed to declare queue"))
	}

	err = ch.ExchangeDeclare("logs", amqp091.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to declare exchange"))
	}

	//deklarasikan bindingnya
	err = ch.QueueBind(q.Name, "", "logs", false, nil)
	if err != nil {
		panic(errors.Wrap(err, "failed to bind queue"))
	}

	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	forever := make(chan struct{})

	go func() {
		for d := range msg {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("[*] waiting for message. To exit press CTRL+C")
	<-forever
}
