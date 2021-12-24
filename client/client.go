package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"

	uc "github.com/LesMiserables1/UFTP-NSQ/usecase"
)

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

	var file []uc.FileTransfer

	for {
		inputBytes := make([]byte, 5*1024)

		n, _, err := connection.ReadFromUDP(inputBytes)
		if err == nil {
			buffer := bytes.NewBuffer(inputBytes[:n])
			dec := gob.NewDecoder(buffer)

			var filePart uc.FileTransfer
			err = dec.Decode(&filePart)

			if err == nil {
				file = append(file, filePart)
			}
			if len(file) == 170 {
				return file
			}
		} else {
			fmt.Println(err)
		}
	}
}
