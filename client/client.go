package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"

	uc "github.com/LesMiserables1/UFTP-NSQ/usecase"
	"github.com/nsqio/go-nsq"
)

type Message struct {
	ArrayMissingParts []int
}

var File []uc.FileTransfer

func receiveMessageUDP() []uc.FileTransfer {

	s, err := net.ResolveUDPAddr("udp", ":3000")
	if err != nil {
		panic(err)
	}
	connection, err := net.ListenUDP("udp", s)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	quit := make(chan struct{})

	for i := 0; i < runtime.NumCPU(); i++ {
		go listen(connection, quit)
	}
	<-quit
	return File
}
func listen(connection *net.UDPConn, quit chan struct{}) {
	inputBytes := make([]byte, 5*1024)
	for File[170].Part == 0 {

		n, _, err := connection.ReadFromUDP(inputBytes)
		if err == nil {
			appendFile(inputBytes, n)
		} else {
			fmt.Println(err)
		}

	}
	quit <- struct{}{}
}
func appendFile(inputBytes []byte, n int) {
	buffer := bytes.NewBuffer(inputBytes[:n])
	dec := gob.NewDecoder(buffer)

	var filePart uc.FileTransfer
	err := dec.Decode(&filePart)

	if err == nil {
		File[filePart.Part] = filePart
	}
}

type myMessageHandler struct{}

func (h *myMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}
	lengthSize, err := strconv.Atoi(string(m.Body))
	if err != nil {
		log.Fatal(err)
	}
	File = make([]uc.FileTransfer, lengthSize)
	sendMessage(lengthSize)
	return nil
}

func sendMessage(lengthFile int) {
	config := nsq.NewConfig()

	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}

	missingPart := Message{
		ArrayMissingParts: findMissingPart(),
	}
	messageBody, err := json.Marshal(missingPart)
	if err != nil {
		log.Fatal(err)
	}
	topicName := "missingFile"
	err = producer.Publish(topicName, messageBody)
	if err != nil {
		log.Fatal(err)
	}
}

func findMissingPart() []int {
	var result []int
	for i, values := range File {
		if values.Part == 0 {
			result = append(result, i)
		}
	}
	return result
}
