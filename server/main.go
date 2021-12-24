package main

import (
	uc "github.com/LesMiserables1/UFTP-NSQ/usecase"
)

func main() {
	const fileName = `Proposal Skripsi - Andre_FinalDraft.pdf`

	fileParts, err := uc.ChunkingFiles(fileName)
	if err != nil {
		panic(err)
	}

	sendMessageUDP(fileParts)
}
