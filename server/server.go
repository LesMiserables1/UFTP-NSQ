package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"strconv"
	"time"

	uc "github.com/LesMiserables1/UFTP-NSQ/usecase"
	nsq "github.com/nsqio/go-nsq"
)

var receiveTime = time.Now()

func sendMessageUDP(fileTransfer []uc.FileTransfer) {

	connection, err := net.Dial("udp", "127.0.0.1:3000")
	if err != nil {
		panic(err)
	}

	defer connection.Close()

	for _, filePart := range fileTransfer {
		var packet bytes.Buffer

		enc := gob.NewEncoder(&packet)

		err := enc.Encode(filePart)
		if err != nil {
			panic(err)
		}

		_, err = connection.Write(packet.Bytes())
		if err != nil {
			panic(err)
		}
	}
}

type myMessageHandler struct{}

func (h *myMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}
	receiveTime = time.Now()
	// do whatever actual message processing is desired

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return nil
}

func sendMessage(lengthFile int) {
	// Instantiate a producer.
	config := nsq.NewConfig()

	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}

	messageBody := []byte(strconv.Itoa(lengthFile))
	topicName := "lengthFile"
	for {
		if time.Since(receiveTime).Seconds() >= 60 {
			err = producer.Publish(topicName, messageBody)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
