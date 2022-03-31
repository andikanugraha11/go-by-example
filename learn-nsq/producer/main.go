package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

type Message struct {
	Name      string
	Content   string
	Timestamp string
}

func main() {
	// nsq initialization
	nsqConf := nsq.NewConfig()

	// create new nsq producer
	producer, err := nsq.NewProducer("127.0.0.1:4150", nsqConf)
	if err != nil {
		log.Fatalf("failed to create nsq producer: %s", err.Error())
	}

	// example new nsq topic
	topic := "my_topic"

	// example nsq message
	msg := Message{
		Name:      "Example message name",
		Content:   "Example message content",
		Timestamp: time.Now().String(),
	}

	// marshaling nsq message into []byte
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("failed to marshaling nsq message: %s", err.Error())
	}

	// publish message
	err = producer.Publish(topic, payload)
	if err != nil {
		log.Fatalf("failed to publish nsq message: %s", err.Error())
	}
}
