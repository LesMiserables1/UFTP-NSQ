package usecase

type FileTransfer struct {
	Part     int
	FileByte []byte
}

const Megabytes = 5000
