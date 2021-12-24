package usecase

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func ChunkingFiles(fileName string) ([]fileTransfer, error) {
	const absPath = `./files/`

	var filePath = absPath + fileName
	dataFile, err := os.Open(filePath)
	if err != nil {
		return []fileTransfer{}, err
	}
	defer dataFile.Close()

	var fileTransfers []fileTransfer

	r := bufio.NewReader(dataFile)

	x := 0
	for {
		buf := make([]byte, megabytes)

		n, err := r.Read(buf)

		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		fileTransfers = append(fileTransfers, fileTransfer{
			Part:     x,
			FileByte: buf,
		})
		x++
	}
	return fileTransfers, nil
}

func MergingFiles(fileTransfers []fileTransfer) {
	fileName := "./files/output.pdf"
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
