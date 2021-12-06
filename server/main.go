package main

import (
	uc "github.com/LesMiserables1/UFTP-NSQ/usecase"
)

func main() {
	const fileName = `Proposal Skripsi - Andre_FinalDraft.pdf`

	_, err := uc.ChunkingFiles(fileName)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%+v", fileParts)
}
