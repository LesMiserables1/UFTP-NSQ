package usecase

type fileTransfer struct {
	Part     int
	FileByte []byte
}

const megabytes = 4 * 1024
