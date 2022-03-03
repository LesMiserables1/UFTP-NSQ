package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"time"

	uc "github.com/LesMiserables1/UFTP-NSQ/usecase"
	"github.com/nsqio/go-nsq"
)

type Message struct {
	ArrayMissingParts []int
}

var FileParts []uc.FileTransfer

func sendMessageUDP(missingPart []int) {

	connection, err := net.Dial("udp", "127.0.0.1:3000")
	if err != nil {
		panic(err)
	}

	defer connection.Close()
	for _, filePart := range missingPart {
		missingFilePart, _ := _map.Load("fileParts")
		missingFilePart = missingFilePart.([]uc.FileTransfer)[filePart]
		var packet bytes.Buffer

		enc := gob.NewEncoder(&packet)

		err := enc.Encode(missingFilePart)
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
		return nil
	}
	_map.LoadOrStore("receiveTime", time.Now())
	if string(m.Body[:]) != "SELESAI" {
		resMessage := Message{}
		err := json.Unmarshal([]byte(m.Body), &resMessage)
		if err != nil {
			panic(err)
		}
		sendMessageUDP(resMessage.ArrayMissingParts)
	}

	return nil
}
func receiveMessage() {
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("missingFile", "channel", config)
	if err != nil {
		log.Fatal(err)
	}
	consumer.AddHandler(&myMessageHandler{})
	err = consumer.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}
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
		receiveTime, _ := _map.Load("receiveTime")
		if time.Since(receiveTime.(time.Time)).Seconds() >= 60 {
			err = producer.Publish(topicName, messageBody)
			if err != nil {
				log.Fatal(err)
			}
		}
		_map.Store("receiveTime", time.Now())
	}
}
