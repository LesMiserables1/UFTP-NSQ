package main

import (
	uc "github.com/LesMiserables1/UFTP-NSQ/usecase"
)

func main() {
	file := receiveMessageUDP()

	uc.MergingFiles(file)
}
