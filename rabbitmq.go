package main

/*
 * @Date: 2020-11-20 20:26:52
 * @LastEditors: kanoyami
 * @LastEditTime: 2020-11-20 23:05:03
 */

import (
	"fmt"

	"github.com/streadway/amqp"
)

// RabbitMqChannel RabbitMqChannel
var RabbitMqChannel *amqp.Channel = nil

func failOnMqError(err error, msg string) {
	if err != nil {
		lPrintErr("%s: %s", msg, err)
	}
}

// SetMqchannel 链接rabbitMQ
func SetMqchannel(name string, pwd string, ip string, port string) *amqp.Channel {
	RabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", name, pwd, ip, port)
	mqConn, err := amqp.Dial(RabbitURL)
	failOnMqError(err, "Failed to connect to RabbitMQ")
	ch, chErr := mqConn.Channel()
	failOnMqError(chErr, "Failed to open a channel")
	lPrintln("RabbitMq连接成功")
	return ch
}

// PublishChannel publicchannel
func PublishChannel(ch *amqp.Channel, info string) {
	q, err := ch.QueueDeclare(
		"aclive_record", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
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
	lPrintln("mqinfo:" + info)
	failOnMqError(err, "Failed to publish a message")
}
