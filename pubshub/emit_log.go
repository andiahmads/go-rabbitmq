package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	for true {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("write your message =")

		mPayload, _ := reader.ReadString('\n')

		conn, err := amqp091.Dial("amqp://userexample:pasexample@portexample:5672/")

		if err != nil {
			panic(errors.Wrap(err, "failed to connect to rabbitmq server"))
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			panic(errors.Wrap(err, "failed to get channel"))
		}
		defer ch.Close()
		// q, err := ch.QueueDeclare("", true, false, true, false, nil)

		// if err != nil {
		// 	panic(errors.Wrap(err, "failed to declare queue"))
		// }
		err = ch.ExchangeDeclare("logs", amqp091.ExchangeFanout, true, false, false, false, nil)
		if err != nil {
			panic(errors.Wrap(err, "failed to declare exchange"))
		}

		err = ch.Publish("logs", "", false, false, amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(mPayload),
		})
	}

}
