package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

const (
	TimeFiveMinute = time.Second * 300
	TimeZeroSecond = time.Second * 0
)

type messageHandler struct{}

type Message struct {
	Name      string
	Content   string
	Timestamp string
}

func main() {
	// nsq initialization
	nsqConf := nsq.NewConfig()

	// nsq configuration
	nsqConf.MaxAttempts = 10
	nsqConf.MaxInFlight = 5
	nsqConf.MaxRequeueDelay = TimeFiveMinute
	nsqConf.DefaultRequeueDelay = TimeZeroSecond

	// nsq topic name listen to
	topic := "my_topic"

	// nsq channel name listen to
	channel := "my_channel"

	// create nsq consumer
	consumer, err := nsq.NewConsumer(topic, channel, nsqConf)
	if err != nil {
		log.Fatalf("failed to create nsq consumer: %s", err.Error())
	}

	// gracefully consumer stop

	// set message handler
	consumer.AddHandler(&messageHandler{})

	// connect to nsqlookup
	consumer.ConnectToNSQLookupd("127.0.0.1:4161")

	// wait exit signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// gracefully consumer stop
	consumer.Stop()
	log.Println("consumer exited")
}

// if message handle returning non-nil error, the cosummer automaticaly send a req command to nsq for re-queue the message based on configuration
func (h *messageHandler) HandleMessage(m *nsq.Message) error {
	var req Message

	// unmarshalling request
	if err := json.Unmarshal(m.Body, &req); err != nil {
		log.Printf("Error marshaling message body: %s", err.Error())
		return err
	}

	// print the message recieved
	log.Println(req)
	return nil
}
