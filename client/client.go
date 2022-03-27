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
	"time"

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
	files, _ := _map.Load("Files")
	return files.([]uc.FileTransfer)
}
func listen(connection *net.UDPConn, quit chan struct{}) {
	inputBytes := make([]byte, 5*1024)
	status, _ := _map.Load("Status")
	for !status.(bool) {
		n, _, err := connection.ReadFromUDP(inputBytes)
		if err == nil {
			_map.Store("receiveTime", time.Now())
			appendFile(inputBytes, n)
		} else {
			fmt.Println(err)
		}
		fmt.Println(status.(bool))
		status, _ = _map.Load("Status")

	}
	quit <- struct{}{}
}
func appendFile(inputBytes []byte, n int) {
	buffer := bytes.NewBuffer(inputBytes[:n])
	dec := gob.NewDecoder(buffer)

	var filePart uc.FileTransfer
	err := dec.Decode(&filePart)

	files, _ := _map.Load("Files")

	if err == nil {
		files.([]uc.FileTransfer)[filePart.Part] = filePart
		_map.Store("Files", files)
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
	_map.Store("Files", File)

	sendMessage(lengthSize)
	return nil
}
func receiveMessage() {
	config := nsq.NewConfig()

	consumer, err := nsq.NewConsumer("lengthFile", "channel", config)
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
	config := nsq.NewConfig()

	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}

	for {
		receiveTime, _ := _map.Load("receiveTime")
		fileMissingArray := findMissingPart()
		missingPart := Message{
			ArrayMissingParts: fileMissingArray,
		}
		topicName := "missingFile"
		if len(missingPart.ArrayMissingParts) == 0 {
			_map.Store("Status", true)
		}
		if time.Since(receiveTime.(time.Time)).Seconds() >= 60 {
			_map.Store("receiveTime", time.Now())
			if len(missingPart.ArrayMissingParts) == 0 {
				messageBody := []byte("SELESAI")
				if err != nil {
					log.Fatal(err)
				}
				err = producer.Publish(topicName, messageBody)
				if err != nil {
					log.Fatal(err)
				}
				_map.Store("Status", true)
			} else {

				messageBody, err := json.Marshal(missingPart)
				if err != nil {
					log.Fatal(err)
				}
				err = producer.Publish(topicName, messageBody)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}

func findMissingPart() []int {
	var result []int
	files, _ := _map.Load("Files")
	for i, values := range files.([]uc.FileTransfer) {
		if len(values.FileByte) == 0 {
			result = append(result, i)
		}

	}
	return result
}
