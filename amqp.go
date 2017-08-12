package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

const (
	AMQPExchange = "events-pub-sub"
)

type AMQP struct {
	exchange string
	ch       *amqp.Channel
}

type Publisher struct {
	AMQP
}

type Subscriber struct {
	AMQP
}

func NewAMQP(userid string) (a AMQP) {
	a.exchange = fmt.Sprintf("%s_%s", AMQPExchange, userid)
	a.Connect()

	return
}

func (a *AMQP) Connect() {
	AMQPConnection, err := amqp.Dial("amqp://fjkjuecr:aOk4Brn5TeB8VpaT861IEGJWuqFPFgbY@golden-kangaroo.rmq.cloudamqp.com/fjkjuecr")
	if err != nil {
		panic(err)
	}

	a.ch, err = AMQPConnection.Channel()
	if err != nil {
		log.Panic(err)
	}

	if err = a.ch.ExchangeDeclare(a.exchange, "fanout", false, true, false, false, nil); err != nil {
		log.Panic(err)
	}

}

func (p Publisher) Publish(m message) error {
	p.Connect()

	if err := p.ch.Confirm(false); err != nil {
		return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
	}

	confirms := p.ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	defer confirmOne(confirms)

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return p.ch.Publish(
		p.exchange, // publish to an exchange
		"",         // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            b,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	)
}

func (s Subscriber) Receive(userid string) (c <-chan amqp.Delivery, err error) {
	s.Connect()

	queue := uuid.NewV4().String()

	_, err = s.ch.QueueDeclare(queue, false, true, true, false, nil)
	if err != nil {
		return
	}

	routingKey := ""
	exchange := fmt.Sprintf("%s_%s", AMQPExchange, userid)

	err = s.ch.QueueBind(queue, routingKey, exchange, false, nil)
	if err != nil {
		return
	}

	return s.ch.Consume(queue, "", false, true, false, false, nil)
}

func confirmOne(confirms <-chan amqp.Confirmation) {
	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
