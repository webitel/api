package broker

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
	api     map[string]chan []byte
	mx      *sync.Mutex
}

type response struct {
	Allow bool `json:"allow"`
}

func (c *Consumer) Send() (r *response) {
	body, _ := json.Marshal(&request{"root", "accounts", "r"})
	uId := uuid.NewV4().String()
	cResp := make(chan []byte, 1)
	c.mx.Lock()
	c.api[uId] = cResp
	c.mx.Unlock()
	c.channel.Publish(
		"engine",       // exchange
		"*.auth.check", // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        []byte(body),
			MessageId:   uId,
			ReplyTo:     "*.auth.check.response",
		})
	<-cResp
	return
}

type request struct {
	Role     string `json:"role"`
	Resource string `json:"resource"`
	Perm     string `json:"perm"`
}

var c *Consumer

func init() {
	c = New()
}

func Send() *response {
	return c.Send()
}

func New() *Consumer {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     "",
		done:    make(chan error),
		api:     make(map[string]chan []byte),
		mx:      &sync.Mutex{},
	}
	var err error

	c.conn, err = amqp.Dial("amqp://webitel:secret@10.10.10.200:5672")
	if err != nil {
		panic(err)
	}

	go func() {
		e := <-c.conn.NotifyClose(make(chan *amqp.Error))
		fmt.Printf("closing: %s", e)
		panic(e)
	}()
	log.Printf("got Connection, getting Channel")

	c.channel, err = c.conn.Channel()
	if err != nil {
		panic(err)
	}

	queue, err := c.channel.QueueDeclare(
		"",    // name of the queue
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)

	if err = c.channel.QueueBind(
		queue.Name,              // name of the queue
		"*.auth.check.response", // bindingKey
		"engine",                // sourceExchange
		true,                    // noWait
		nil,                     // arguments
	); err != nil {
		panic(err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		"",         // consumerTag,
		true,       // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		panic(err)
	}
	if deliveries != nil {

	}

	go handle(deliveries, c.done, c)
	return c
}

func handle(deliveries <-chan amqp.Delivery, done chan error, c *Consumer) {
	for d := range deliveries {
		c.mx.Lock()
		if ch, ok := c.api[d.MessageId]; ok {
			//mx.Lock()
			ch <- d.Body
			//mx.Unlock()
			delete(c.api, d.MessageId)
		}
		c.mx.Unlock()
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
