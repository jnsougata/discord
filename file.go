package discord

import (
	"os"
)

type File struct {
	Name        string
	Content     []byte
	Description string
}

func (f *File) Write(path string) error {
	content, err := os.ReadFile(path)
	f.Content = content
	return err
}
