package usecase

type FileTransfer struct {
	Part     int
	FileByte []byte
}

const Megabytes = 4 * 1024
