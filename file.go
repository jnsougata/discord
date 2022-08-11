package discord

import (
	"fmt"
	"log"
	"os"
)

type File struct {
	Name        string
	Content     []byte
	Description string
}

func (f *File) AddContent(path string) {
	bs, err := os.ReadFile(path)
	if err != nil {
		log.Println(fmt.Sprintf(`Error reading file for path %v`, path))
	}
	f.Content = bs
}
