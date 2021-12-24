package main

import (
	"bytes"
	"encoding/gob"
	"net"
	"time"

	uc "github.com/LesMiserables1/UFTP-NSQ/usecase"
)

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
		time.Sleep(10 * time.Millisecond)

	}

}
