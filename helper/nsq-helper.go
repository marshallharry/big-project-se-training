package helper

import (
	"log"

	"github.com/bitly/go-nsq"
)

var config *nsq.Config

func init() {
	config = nsq.NewConfig()
}

// Publish handler
func Publish(topic string, message []byte) {
	p, _ := nsq.NewProducer("devel-go.tkpd:4150", config)

	err := p.Publish(topic, message)
	if err != nil {
		log.Print("Could not connect")
	}

	p.Stop()
}

// Subscribe handler
func Subscribe(topic string, channel string, handler nsq.HandlerFunc) {
	q, _ := nsq.NewConsumer(topic, channel, config)
	q.AddHandler(handler)
	err := q.ConnectToNSQD("devel-go.tkpd:4150")
	if err != nil {
		log.Print("Could not connect")
	}
}
