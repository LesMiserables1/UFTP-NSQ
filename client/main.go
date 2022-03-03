package main

import (
	"sync"

	uc "github.com/LesMiserables1/UFTP-NSQ/usecase"
)

var _map sync.Map

func main() {

	go receiveMessage()

	file := receiveMessageUDP()
	uc.MergingFiles(file)
}
