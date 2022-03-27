package main

import (
	"sync"
	"time"

	uc "github.com/LesMiserables1/UFTP-NSQ/usecase"
)

var _map sync.Map

func main() {
	var receiveTime = time.Now().Add(time.Second * -60)
	_map.Store("receiveTime", receiveTime)
	_map.Store("Status", false)
	go receiveMessage()

	file := receiveMessageUDP()
	uc.MergingFiles(file)
}
