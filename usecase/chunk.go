package usecase

import (
	"fmt"
	"io"
	"os"
)

func ChunkingFiles(fileName string) ([]FileTransfer, error) {
	const absPath = `../files/`
	var filePath = absPath + fileName
	dataFile, err := os.Open(filePath)
	if err != nil {
		return []FileTransfer{}, err
	}
	defer dataFile.Close()

	var fileTransfers []FileTransfer

	x := 0
	for {
		buf := make([]byte, PartSize)
		n, err := dataFile.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		fileTransfers = append(fileTransfers, FileTransfer{
			Part:     x,
			FileByte: buf,
		})
		x++
	}
	return fileTransfers, nil
}

func MergingFiles(fileTransfers []FileTransfer) {
	fileName := "../files/output.pdf"
	_, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	for _, values := range fileTransfers {
		n, err := file.Write(values.FileByte)

		if err != nil {
			panic(err)
		}
		file.Sync()
		fmt.Println("Written ", n, " bytes")

	}
	defer file.Close()
}
