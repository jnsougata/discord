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

func Content(path string) []byte {
	// Opens the file from the given path
	// Returns the content of the file
	// if the file does not exist, returns empty byte array
	bs, err := os.ReadFile(path)
	if err != nil {
		log.Println(fmt.Sprintf(`Error reading file for path %v`, path))
		return []byte{}
	}
	return bs
}
