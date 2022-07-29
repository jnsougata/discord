package file

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

func Read(path string) []byte {
	bs, err := os.ReadFile(path)
	if err != nil {
		log.Println(fmt.Sprintf(`Error reading file for path %v`, path))
		return []byte{}
	}
	return bs
}
