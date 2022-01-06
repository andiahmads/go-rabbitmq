package main

import (
	"bytes"
	"log"
	"time"

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
	q, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)

	if err != nil {
		panic(errors.Wrap(err, "failed to declare queue"))
	}
	//jumlah pekerjaan queue
	err = ch.Qos(1, 0, false)
	if err != nil {
		panic(errors.Wrap(err, "failed to set QOS"))

	}

	msg, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	forever := make(chan struct{})

	go func() {
		for d := range msg {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Println("Done!")
			d.Ack(false)
		}
	}()

	log.Printf("[*] waiting for message. To exit press CTRL+C")
	<-forever
}
