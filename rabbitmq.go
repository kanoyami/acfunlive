package main

/*
 * @Date: 2020-11-20 20:26:52
 * @LastEditors: kanoyami
 * @LastEditTime: 2020-11-20 21:07:58
 */

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var RabbitMqConn = nil

func failOnMqError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// SetMqchannel 链接rabbitMQ
func SetMqchannel(name string, pwd string, ip string, port string) *amqp.Channel {
	RabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", name, pwd, ip, port)
	mqConn, err := amqp.Dial(RabbitURL)
	failOnMqError(err, "Failed to connect to RabbitMQ")
	defer mqConn.Close()
	ch, chErr := mqConn.Channel()
	failOnMqError(chErr, "Failed to open a channel")
	defer ch.Close()
	return ch
}

// PublishChannel publicchannel
func PublishChannel(ch *amqp.Channel, info string) {
	q, err := ch.QueueDeclare(
		"aclive_record_final", // name
		false,                 // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	failOnMqError(err, "Failed to declare a queue")

	body := info
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnMqError(err, "Failed to publish a message")
}
