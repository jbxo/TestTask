package main

import (
	"os"
)

type FileSaver struct {
	file *os.File
}

func NewFileSaver(filePath string) *FileSaver {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	return &FileSaver{file}
}

func (fs *FileSaver) Close() error {
	return fs.file.Close()
}

func (fs *FileSaver) SaveData(data []byte) error {
	_, err := fs.file.Write(data)
	return err
}
