package file

import (
	"bufio"
	"os"
)

type File struct {
	*os.File
	*bufio.Writer
	*bufio.Reader
}

func NewFile() *File {
	return &File{}
}

func (e *File) Open(filePath string) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	e.File = file
	e.Writer = bufio.NewWriter(file)
	e.Reader = bufio.NewReader(file)
}

func (e *File) Close() {
	if e.File != nil {
		e.File.Close()
	}
}
