package usecase

type FileTransfer struct {
	Part     int
	FileByte []byte
}

const PartSize = 5000
