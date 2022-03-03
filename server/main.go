package main

import (
	"runtime"
	"sync"
	"time"

	uc "github.com/LesMiserables1/UFTP-NSQ/usecase"
)

var _map sync.Map

func main() {

	const fileName = `Proposal Skripsi - Andre_FinalDraft.pdf`

	FileParts, err := uc.ChunkingFiles(fileName)
	if err != nil {
		panic(err)
	}
	var receiveTime = time.Now().Add(time.Second * -60)

	_map.Store("fileParts", FileParts)
	_map.Store("receiveTime", receiveTime)
	go sendMessage(len(FileParts))
	receiveMessage()

	runtime.Goexit()

}
